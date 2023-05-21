package http

import (
	entity "bili_danmaku/internal/types"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"strings"
	"time"
)

func GetLoginUrl() (*entity.LoginUrl, error) {
	var err error
	var resp *resty.Response
	var url = "https://passport.bilibili.com/qrcode/getLoginUrl"

	r := &entity.LoginUrl{}
	if resp, err = cli.R().
		SetHeader("user-agent", userAgent).
		Get(url); err != nil {
		logx.Error("请求getLoginUrl失败：", err)
		return nil, err
	}
	if err = json.Unmarshal(resp.Body(), r); err != nil {
		logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
		return nil, err
	}

	logx.Info("oauthKey:", r.Data.OauthKey)

	return r, err
}

// 验证登录的同时，将cookie赋值
func GetLoginInfo(oauthKey string) (*entity.LoginInfoData, error) {
	var err error
	var url = "https://passport.bilibili.com/qrcode/getLoginInfo?oauthKey=" + oauthKey
	var resp *resty.Response
	var data *entity.LoginInfoData
	var file *os.File
	pre := &entity.LoginInfoPre{}
	logx.Info("等待扫码登录...")
	for {
		logx.Info(1)
		if resp, err = cli.R().
			SetHeader("user-agent", userAgent).
			Post(url); err != nil {
			logx.Error("请求getLoginInfo失败：", err)
			return nil, err
		}

		if err = json.Unmarshal(resp.Body(), pre); err != nil {
			logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
			return nil, err
		}

		if pre.Status {
			data = &entity.LoginInfoData{}
			if err = json.Unmarshal(resp.Body(), data); err != nil {
				logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
				return nil, err
			}
			logx.Info("登录成功！")
			break
		}

		time.Sleep(5 * time.Second)
	}

	for _, v := range resp.Header().Values("Set-Cookie") {
		pair := strings.Split(v, ";")
		kv := strings.Split(pair[0], "=")
		CookieList[kv[0]] = kv[1]
		CookieStr += pair[0] + ";"
	}
	//使用追加模式打开文件
	file, err = os.OpenFile("token/bili_token.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logx.Errorf("打开文件错误：", err)
	}
	file.WriteString(CookieStr)
	file.Close()
	//使用追加模式打开文件
	file, err = os.OpenFile("token/bili_token.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logx.Errorf("打开文件错误：", err)
	}
	tokenstr, _ := json.Marshal(CookieList)
	file.WriteString(string(tokenstr))
	file.Close()
	return data, err
}

func FileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
func SetHistoryCookie() error {
	var cookie []byte
	var err error
	cookie, err = os.ReadFile("token/bili_token.txt")
	if err != nil {
		logx.Errorf("打开文件错误：", err)
		return err
	}
	CookieStr = string(cookie)
	cookie, err = os.ReadFile("token/bili_token.json")
	if err != nil {
		logx.Errorf("打开文件错误：", err)
		return err
	}
	err = json.Unmarshal(cookie, &CookieList)
	if err != nil {

		return err
	}
	//CookieStr = string(cookie)
	return err
}
