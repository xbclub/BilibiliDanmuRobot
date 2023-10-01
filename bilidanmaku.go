package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"os"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed etc/bilidanmaku-api.yaml
var configfileembed embed.FS
var configFile = flag.String("f", "etc/bilidanmaku-api.yaml", "the config file")
var Version string

func init() {

}
func main() {
	flag.Parse()
	var c config.Config
	dir := "./token"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.Mkdir(dir, 0755)
		if err != nil {
			panic(fmt.Sprintf("无法创建token文件夹 请手动创建:%s", err))
		}
	}
	dir = "./etc"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.Mkdir(dir, 0755)
		if err != nil {
			panic(fmt.Sprintf("无法创建token文件夹 请手动创建etc文件夹:%s", err))
		}
	}
	exists, err := fileExists(dir + "/bilidanmaku-api.yaml")
	if err != nil {
		// 处理错误
		panic(err)
	}
	if exists == false {
		file, err := configfileembed.Open("etc/bilidanmaku-api.yaml")
		if err != nil {
			panic(err)
		}
		destFile, err := os.Create(dir + "/bilidanmaku-api.yaml")
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(destFile, file)
		if err != nil {
			panic(err)
		}
		destFile.Close()
		file.Close()
	}
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	logx.MustSetup(c.Log)
	logx.DisableStat()
	ctx := svc.NewServiceContext(c)
	// Create an instance of the app structure
	logx.Info("当前版本:", Version)
	app := NewApp(Version)

	program := NewProgram(ctx)
	//test := NewTest()
	// Create application with options
	err = wails.Run(&options.App{
		Title:  "弹幕机器人",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		//BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: app.startup,
		//OnBeforeClose: app.Onstop,
		Bind: []interface{}{
			app,
			//test,
			program,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
func fileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		// 文件存在
		return true, nil
	}
	if os.IsNotExist(err) {
		// 文件不存在
		return false, nil
	}
	// 其他错误，比如权限问题等
	return false, err
}
