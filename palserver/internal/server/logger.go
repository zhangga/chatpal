package server

import (
	"fmt"
	"github.com/zhangga/chatpal/palserver/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

func initLogger(conf *config.ConfLog) (*zap.Logger, error) {
	zapLoggerEncoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller",
		MessageKey:       "message",
		StacktraceKey:    "stacktrace",
		EncodeCaller:     customCallerEncoder,
		EncodeTime:       customTimeEncoder,
		EncodeLevel:      customLevelEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		LineEnding:       "\n",
		ConsoleSeparator: " ",
	}

	// 日志等级
	level, err := zapcore.ParseLevel(conf.Level)
	if err != nil {
		return nil, err
	}

	// 日志编码
	var encoder zapcore.Encoder
	if conf.IsJsonEncoder {
		encoder = zapcore.NewJSONEncoder(zapLoggerEncoderConfig)
	} else {
		zapLoggerEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(zapLoggerEncoderConfig)
	}

	var logWriters []zapcore.WriteSyncer
	if conf.IsShowInConsole {
		logWriters = append(logWriters, zapcore.AddSync(os.Stdout))
	}
	logWriters = append(logWriters, &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:  fmt.Sprintf("logs/%s.log", conf.FileName), // ⽇志⽂件路径
			MaxSize:   conf.MaxSize,                              // 单位为MB,默认为512MB
			MaxAge:    conf.MaxAge,                               // 文件最多保存多少天
			LocalTime: true,                                      // 采用本地时间
			Compress:  false,                                     // 是否压缩日志
		}),
		Size: 4096,
	})

	zapCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(logWriters...), level)
	zapLogger := zap.New(zapCore, zap.AddCaller(), zap.AddCallerSkip(0))
	return zapLogger, nil
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func customLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(l.CapitalString())
}

func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}
