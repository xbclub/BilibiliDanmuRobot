package main

import (
	"bili_danmaku/internal/config"
	"bili_danmaku/internal/logic"
	"bili_danmaku/internal/svc"
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

var configFile = flag.String("f", "etc/bilidanmaku-api.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	logx.MustSetup(c.Log)
	logx.DisableStat()
	dir := "./token"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.Mkdir(dir, 0755)
		if err != nil {
			panic(fmt.Sprintf("无法创建token文件夹 请手动创建:%s", err))
			return
		}
	}
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
