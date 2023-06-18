package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

//go:embed all:frontend/dist
var assets embed.FS
var configFile = flag.String("f", "etc/bilidanmaku-api.yaml", "the config file")
var Version string

func init() {
	dir := "./token"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.Mkdir(dir, 0755)
		if err != nil {
			panic(fmt.Sprintf("无法创建token文件夹 请手动创建:%s", err))
		}
	}
	logconf := logx.LogConf{
		ServiceName:         "",
		Mode:                "file",
		Encoding:            "plain",
		Path:                "./logs/applog",
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
	logx.MustSetup(logconf)
	logx.DisableStat()
}
func main() {
	// Create an instance of the app structure
	logx.Info("当前版本:", Version)
	app := NewApp()
	program := NewProgram()
	//test := NewTest()
	// Create application with options
	err := wails.Run(&options.App{
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
