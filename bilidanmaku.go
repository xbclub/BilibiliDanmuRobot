package main

import (
	"bili_danmaku/internal/config"
	"bili_danmaku/internal/logic"
	"bili_danmaku/internal/svc"
	"context"
	"embed"
	"flag"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/zeromicro/go-zero/core/conf"
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
			return
		}
	}
}
func start() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	logx.MustSetup(c.Log)
	logx.DisableStat()
	logx.Infof("当前版本: %s", Version)
	logx.Infof("监听直播间: %d", c.RoomId)

	ctx := svc.NewServiceContext(c)
	server := logic.NewBili_danmakuLogic(context.TODO(), ctx)
	server.Bili_danmaku_Start()

	//server := rest.MustNewServer(c.RestConf)
	defer server.Bili_danmaku_Stop()

	//handler.RegisterHandlers(server, ctx)
	//logx.Infof()
	//fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	//server.Start()
}
func main() {
	// Create an instance of the app structure
	app := NewApp()
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
		Bind: []interface{}{
			app,
			//test,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
