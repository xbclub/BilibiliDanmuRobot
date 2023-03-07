package bullet_girl

import (
	"context"
	"sync"
	"time"
)

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
	var w = 1 * time.Second
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
	for _, v := range interractGiver.handlermsg {
		if len(v) > 13 {
			PushToBulletSender("[欢迎 " + v[0:6] + " ~]")
		} else {
			PushToBulletSender("[欢迎 " + v[0:6] + " ~]")
		}
	}
}
