package main

import (
	"bili_danmaku/internal/config"
	"bili_danmaku/internal/logic"
	"bili_danmaku/internal/svc"
	"context"
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/bilidanmaku-api.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	logx.MustSetup(c.Log)
	logx.DisableStat()
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
