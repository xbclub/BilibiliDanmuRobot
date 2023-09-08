package http

import (
	"github.com/go-resty/resty/v2"
)

const (
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36 Edg/116.0.1938.76"
	userAgent = UserAgent
)

// 调用GetLoginInfo后，对全局变量cookie赋值
var CookieStr string
var CookieList = make(map[string]string)

// 全局客户端对象
var cli *resty.Client

func InitHttpClient() {
	cli = resty.New()
}
