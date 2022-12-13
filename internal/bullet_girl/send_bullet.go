package bullet_girl

import (
	"bili_danmaku/internal/http"
	"bili_danmaku/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

var sender *BulletSender

type BulletSender struct {
	bulletChan chan string
}

func PushToBulletSender(bullet string) {
	logx.Info("PushToBulletSender成功", bullet)
	sender.bulletChan <- bullet
}

func StartSendBullet(ctx context.Context, svcCtx *svc.ServiceContext) {
	var err error

	sender = &BulletSender{
		bulletChan: make(chan string, 1000),
	}

	var msg string
	for {
		select {
		case <-ctx.Done():
			goto END
		case msg = <-sender.bulletChan:
			var danmuLen = svcCtx.Config.DanmuLen
			var msgdata []string
			msgrun := []rune(msg)
			msgLen := len(msgrun)
			msgcount := msgLen / danmuLen
			tmpmsgcount := msgLen % danmuLen
			if tmpmsgcount != 0 {
				msgcount += 1
			}
			for m := 1; m <= msgcount; m++ {
				if msgLen < m*danmuLen {
					msgdata = append(msgdata, string(msgrun[(m-1)*danmuLen:msgLen]))
					continue
				}
				msgdata = append(msgdata, string(msgrun[(m-1)*danmuLen:danmuLen*m]))
			}
			for _, msgs := range msgdata {
				if err = http.Send(msgs, svcCtx); err != nil {
					logx.Errorf("弹幕发送失败：", err, "msg:", msg)
				} else {
					logx.Infof("弹幕发送成功：", msg)
				}
				//fmt.Println(msgs)
				time.Sleep(1 * time.Second) // 防止弹幕发送过快
			}
		}

	}
END:
}
