package room

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

var (
	AppId     string
	AppSecret string
)

func sendHttpRequest(method, url string, body any) ([]byte, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 读取响应内容
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func sendHttpGet(baseUrl string, param any) ([]byte, error) {
	// 构建查询字符串
	queryParams, err := structToURLValues(param)
	// 完整的URL
	fullUrl := fmt.Sprintf("%s?%s", baseUrl, queryParams.Encode())
	// 发送GET请求
	resp, err := http.Get(fullUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

// structToURLValues 将结构体转换为URL查询参数
func structToURLValues(data interface{}) (url.Values, error) {
	values := url.Values{}
	v := reflect.ValueOf(data)

	// 遍历结构体的所有字段
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typeField := v.Type().Field(i)
		tag := typeField.Tag.Get("url")

		// 如果字段有'url'标签，则作为参数名
		if tag != "" && field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface() {
			values.Set(tag, fmt.Sprintf("%v", field.Interface()))
		}
	}

	return values, nil
}
