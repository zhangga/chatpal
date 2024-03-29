package room

import (
	"github.com/goccy/go-json"
	"github.com/zhangga/chatpal/palserver/internal/msg"
)

type LoginReq struct {
	Code string `json:"code,omitempty"`
}
type LoginResp struct {
	Code int `json:"code,omitempty"`
}

// handlerLogin 登录处理
func handleLogin(c *conn, msgType int, content []byte) (msg.Message, error) {
	var req LoginReq
	if err := json.Unmarshal(content, &req); err != nil {
		c.logger.Errorf("error while unmarshaling login request, err: %v", err)
		return LoginResp{Code: 1}, err
	}

	type wxLoginReq struct {
		AppId     string `url:"appid"`
		Secret    string `url:"secret"`
		JsCode    string `url:"js_code"`
		GrantType string `url:"grant_type"`
	}
	wxReq := wxLoginReq{
		AppId:     AppId,
		Secret:    AppSecret,
		JsCode:    req.Code,
		GrantType: "authorization_code",
	}
	body, err := sendHttpGet("https://api.weixin.qq.com/sns/jscode2session", wxReq)
	if err != nil {
		c.logger.Errorf("error while sending http request, err: %v", err)
		return LoginResp{Code: 2}, err
	}

	type wxLoginResp struct {
		SessionKey string `json:"session_key,omitempty"`
		OpenId     string `json:"openid,omitempty"`
		Unionid    string `json:"unionid,omitempty"`
		ErrMsg     string `json:"errmsg,omitempty"`
		ErrCode    int32  `json:"errcode,omitempty"`
	}
	var wxResp wxLoginResp
	if err = json.Unmarshal(body, &wxResp); err != nil {
		c.logger.Errorf("error while unmarshaling wx login response, err: %v", err)
		return LoginResp{Code: 3}, err
	}

	c.logger.Debugf("wx login response: %#v", wxResp)
	user := &User{
		ConnId: c.id,
		OpenId: wxResp.OpenId,
	}
	c.user = user
	return LoginResp{Code: 0}, nil
}
