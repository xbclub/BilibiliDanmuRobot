package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/http"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/yaml.v3"
	"io"
	httpx "net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// App struct
type App struct {
	Vserion      string
	ctx          context.Context
	login        chan bool
	loginurl     *entity.LoginUrl
	loginstatus  int
	channelisrun bool
	loginCtx     context.Context
	loginCancel  context.CancelFunc
}

// NewApp creates a new App application struct
func NewApp(version string) *App {
	http.InitHttpClient()
	return &App{Vserion: version}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
func (l *App) Userlogin() *entity.LoginUrl {
	if l.GetloginStatus() {
		return nil
	}
	var err error
	if l.login != nil && l.channelisrun {
		l.Stopwork()
	}
	l.login = make(chan bool)
	if l.loginurl, err = http.GetLoginUrl(); err != nil {
		logx.Error(err)
		return nil
	}
	l.loginstatus = 0
	l.loginCtx, l.loginCancel = context.WithCancel(context.Background())
	go l.work(l.loginCtx)
	l.channelisrun = true
	return l.loginurl
	//if err = utiles.GenerateQr(loginUrl.Data.Url); err != nil {
	//	logx.Error(err)
	//	return err
	//}
}
func (l *App) Getlogin() int {
	if l.loginstatus == 1 {
		l.Stopwork()
		return 1
	} else if l.loginstatus == 3 {
		l.Stopwork()
		return 3
	}
	return 2
}
func (l *App) work(ctx context.Context) {
	var err error
	var url = "https://passport.bilibili.com/qrcode/getLoginInfo?oauthKey=" + l.loginurl.Data.OauthKey
	var resp *resty.Response
	var data *entity.LoginInfoData
	var file *os.File
	var CookieStr string
	var CookieList = make(map[string]string)
	cli := resty.New()
	pre := &entity.LoginInfoPre{}
	logx.Info("等待扫码登录...")
	nologin := true
	var w = 1 * time.Second
	var t = time.NewTimer(w)
	defer t.Stop()
	for nologin {
		select {
		case <-ctx.Done():
			//l.channelisrun = false
			logx.Info("登录线程退出")
			nologin = false
		case <-t.C:
			logx.Info("等待扫码中")

			if resp, err = cli.R().
				SetHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.64 Safari/537.36").
				Post(url); err != nil {
				logx.Error("请求getLoginInfo失败：", err)
				l.loginstatus = 3
			}

			if err = json.Unmarshal(resp.Body(), pre); err != nil {
				logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
				l.loginstatus = 3
			}

			if pre.Status {

				data = &entity.LoginInfoData{}
				if err = json.Unmarshal(resp.Body(), data); err != nil {
					logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
					l.loginstatus = 3
				}
				logx.Info("登录成功！")
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
					l.loginstatus = 3
				}
				tokenstr, _ := json.Marshal(CookieList)
				file.WriteString(string(tokenstr))
				file.Close()
				l.loginstatus = 1
				l.Stopwork()
				//l.Stopwork()
			}
			t.Reset(w)
		}
	}
}
func (l *App) Stopwork() {
	if l.loginCancel != nil {
		l.loginCancel()
	}
}

func (l *App) Onstop() bool {
	return true
}
func (l *App) GetloginStatus() bool {
	if http.FileExists("token/bili_token.txt") && http.FileExists("token/bili_token.json") {
		err := http.SetHistoryCookie()
		if err != nil {
			logx.Error(err)
			return false
		}
		status := http.GetUserInfo()
		return status.Islogin
	} else {
		return false
	}

}
func (l *App) GetUserInfo() *entity.UserinfoLite {
	return http.GetUserInfo()
}
func (l *App) WriteConfig(data string) *ConfigResponse {
	var c config.Config
	resp := new(ConfigResponse)
	err := json.Unmarshal([]byte(data), &c)
	if err != nil {
		logx.Error(err)
		resp.Code = false
		resp.Msg = err.Error()
		return resp
	}
	c.Log = logx.LogConf{
		ServiceName:         "",
		Mode:                "file",
		Encoding:            "plain",
		Path:                fmt.Sprintf("./logs/%v", c.RoomId),
		TimeFormat:          "",
		Level:               "info",
		MaxContentLength:    0,
		Compress:            false,
		Stat:                true,
		KeepDays:            0,
		StackCooldownMillis: 100,
		MaxBackups:          0,
		MaxSize:             0,
		Rotation:            "daily",
	}
	yamlBytes, err := yaml.Marshal(&c)
	if err != nil {
		logx.Error("Failed to marshal YAML: %v", err)
		resp.Code = false
		resp.Msg = err.Error()
		return resp
	}
	if _, err = os.Stat("./etc"); os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.Mkdir("./etc", 0755)
		if err != nil {
			logx.Error(err)
			resp.Code = false
			resp.Msg = err.Error()
			return resp
		}
	}
	file, err := os.OpenFile("etc/bilidanmaku-api.yaml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logx.Errorf("打开文件错误：", err)
		resp.Code = false
		resp.Msg = "打开文件错误：" + err.Error()
		return resp
	}
	_, err = file.Write(yamlBytes)
	if err != nil {
		logx.Errorf("文件写入错误：", err)
		resp.Code = false
		resp.Msg = "文件写入错误：" + err.Error()
		return resp
	}
	file.Close()

	err = Mustload(&c)
	if err != nil {
		logx.Error(err)
		resp.Code = false
		resp.Msg = err.Error()
		return resp
	}
	resp.Code = true
	return resp
}
func Mustload(c *config.Config) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// 处理panic错误
			// 如果需要，将错误转换为适当的错误类型
			// 设置错误消息或执行任何必要的错误处理
			// ...

			// 返回错误
			// 可以返回自定义的错误消息，或者根据捕获到的panic错误创建一个新的错误对象
			errMsg := fmt.Sprintf("发生了panic错误：%v", r)
			// 或者创建一个新的错误对象：err = errors.New("发生了panic错误")

			// 返回错误
			// someFunction的调用者将接收到该错误
			err = errors.New(errMsg)
		}
	}()
	conf.MustLoad("etc/bilidanmaku-api.yaml", c, conf.UseEnv())
	return nil
}
func (l *App) ReadConfig() *ConfigResponse {
	resp := new(ConfigResponse)
	var c config.Config
	err := Mustload(&c)
	if err != nil {
		logx.Error(err)
		resp.Code = false
		resp.Msg = err.Error()
		return resp
	}
	resp.Code = true
	resp.Form = c
	return resp
}

type UpdateResponse struct {
	Version   string `json:"version"`
	Link      string `json:"link"`
	ChangeLog string `json:"changeLog"`
}

func (l *App) GetVersion() string {
	return l.Vserion
}
func (l *App) CheckUpdate() *VersionResponse {
	versionresp := VersionResponse{Code: 2}

	resp, err := httpx.Get("https://danmuji.neuedu.work/getUpdate")
	if err != nil {
		logx.Error(err)
		versionresp.Msg = "链接更新服务器失败"
		return &versionresp
	}
	if resp.StatusCode != httpx.StatusOK {
		logx.Error(resp.StatusCode)
		versionresp.Msg = "链接更新服务器失败"
		return &versionresp
	}
	updateResp := &UpdateResponse{}
	err = json.NewDecoder(resp.Body).Decode(updateResp)
	if err != nil {
		logx.Error(err)
		versionresp.Msg = "更信息读取失败"
		return &versionresp
	}
	if l.Vserion == updateResp.Version {
		versionresp.Code = 0
		versionresp.Msg = "无需更新"
		return &versionresp
	}
	versionresp.Code = 1
	versionresp.Msg = updateResp.Version
	versionresp.Content = updateResp.ChangeLog
	return &versionresp
}
func (l *App) GetUpdateUpgrader() string {
	resp, err := httpx.Get("https://danmuji.neuedu.work/getUpgraderUpdate")
	if err != nil {
		logx.Error(err)
		return "链接更新服务器失败"
	}
	defer resp.Body.Close()
	if resp.StatusCode == httpx.StatusOK {
		updateResp := &UpdateResponse{}
		err := json.NewDecoder(resp.Body).Decode(updateResp)
		if err != nil {
			logx.Error(err)
			return "更信息读取失败"
		}
		err = downloadAndExtract(updateResp.Link)
		if err != nil {
			logx.Error("Error:", err)
			return "下载更新程序失败"
		}
	} else {
		logx.Errorf("Request failed with status code: %d\n", resp.StatusCode)
		return "更新服务器链接失败"
	}
	go func() {
		// 唤起upgrade.exe start是非阻塞 /B是隐藏窗口
		cmd := exec.Command("cmd.exe", "/C", "start", "upgrader.exe")
		if err := cmd.Start(); err != nil {
			fmt.Println(err)
		}
		fmt.Println("exit")
		os.Exit(0)
	}()
	return "升级程序初始化成功，即将开始升级"
}

//	func (l *App) StartProgram() bool {
//		err := l.program.Start()
//		if err != nil {
//			return false
//		}
//		return true
//	}
//
//	func (l *App) GetProgramStatus() bool {
//		return l.GetProgramStatus()
//	}
type ConfigResponse struct {
	Code bool
	Msg  string
	Form config.Config
}
type VersionResponse struct {
	Code    int
	Msg     string
	Content string
}

func downloadAndExtract(link string) error {
	resp, err := httpx.Get(link)
	logx.Info(link)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	extractedFile, err := os.OpenFile(filepath.Base("upgrader.exe"), os.O_WRONLY|os.O_CREATE, 0755)
	defer extractedFile.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(extractedFile, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
