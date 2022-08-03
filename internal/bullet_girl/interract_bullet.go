package bullet_girl

import (
	"context"
	"sync"
	"time"
)

// 检测到礼物，push [uname]->[giftName]->[cost]，number+1
// 每10s统计一次礼物，并进行感谢，礼物价值高于x元加一句大气

var interractGiver *InterractGiver

type InterractGiver struct {
	tmpmsg     []string
	handlermsg []string
	tableMu    sync.RWMutex
	giftChan   chan *string
}

func pushToInterractChan(g *string) {
	interractGiver.giftChan <- g
}

func Interact(ctx context.Context) {

	interractGiver = &InterractGiver{
		tmpmsg:     []string{},
		handlermsg: []string{},
		tableMu:    sync.RWMutex{},
		giftChan:   make(chan *string, 1000),
	}

	var g *string
	var w = 10 * time.Second
	var t = time.NewTimer(w)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			goto END
		case <-t.C:
			interractGiver.handlermsg = interractGiver.tmpmsg
			interractGiver.tmpmsg = []string{}
			//if rand.Intn(100) < 30 {
			handleInterract()
			//}
			interractGiver.handlermsg = []string{}
			t.Reset(w)
		case g = <-interractGiver.giftChan:
			interractGiver.tmpmsg = append(interractGiver.tmpmsg, *g)
		}
	}
END:
}

func handleInterract() {
	msg := ""
	for k, v := range interractGiver.handlermsg {
		if k == 0 {
			msg += "欢迎 " + v
		} else {
			msg += "，" + v
		}

	}
	if len(msg) == 0 {
		return
	}
	PushToBulletSender(msg + " 进入直播间")
}
