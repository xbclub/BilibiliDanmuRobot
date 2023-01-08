package http

import (
	entity "bili_danmaku/internal/types"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"net/url"
	"time"
	"unicode/utf8"
)

// 调用青云客机器人api
func RequestQingyunkeRobot(msg string) (string, error) {
	var err error
	var urls = "http://api.qingyunke.com/api.php?key=free&appid=0&msg=" + encodeSpecialChar(msg) + "&_=" + fmt.Sprint(time.Now().UnixMicro())
	var resp *resty.Response

	if resp, err = cli.R().
		SetHeader("Content-Type", "utf-8").
		Get(urls); err != nil {
		logx.Error("请求qingyunke机器人接口失败：", err)
		return "", err
	}

	r := &entity.QinugyunkeRobotReplay{}
	err = json.Unmarshal(resp.Body(), r)

	return r.Content, err
}
func encodeSpecialChar(s string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			buffer.WriteRune(r)
		} else if r < utf8.RuneSelf {
			buffer.WriteString(url.QueryEscape(url.QueryEscape(string(r))))
		} else {
			buffer.WriteString(string(r))
		}
		i += size
	}
	return buffer.String()
}
