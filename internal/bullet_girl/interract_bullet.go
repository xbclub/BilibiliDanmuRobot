package bullet_girl

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

var interractGiver *InterractGiver

type InterractGiver struct {
	interractFilter map[int64]time.Time
	locked          *sync.Mutex
	tableMu         sync.RWMutex
	interractChan   chan *InterractData
}
type InterractData struct {
	Uid int64
	Msg string
}

func pushToInterractChan(g *InterractData) {
	interractGiver.interractChan <- g
}

func Interact(ctx context.Context) {

	interractGiver = &InterractGiver{
		interractFilter: map[int64]time.Time{},
		locked:          new(sync.Mutex),
		tableMu:         sync.RWMutex{},
		interractChan:   make(chan *InterractData, 1000),
	}

	var g *InterractData
	var w = 10 * time.Second
	var t = time.NewTimer(w)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			goto END
		case <-t.C:
			//interractGiver.handlermsg = interractGiver.tmpmsg
			//interractGiver.tmpmsg = []string{}
			////if rand.Intn(100) < 30 {
			//handleInterract()
			////}
			//interractGiver.handlermsg = []string{}
			if len(interractGiver.interractFilter) > 0 {
				interractGiver.locked.Lock()
				for k, v := range interractGiver.interractFilter {
					if v.Add(w).Unix() < time.Now().Unix() {
						delete(interractGiver.interractFilter, k)
						logx.Debugf("用户 %v 已从重复过滤列表移除", k)
					}
				}
				interractGiver.locked.Unlock()
			}

			t.Reset(w)
		case g = <-interractGiver.interractChan:
			//interractGiver.tmpmsg = append(interractGiver.tmpmsg, *g)
			interractGiver.locked.Lock()
			if value, ok := interractGiver.interractFilter[g.Uid]; ok && value.Add(w).Unix() >= time.Now().Unix() {
				logx.Debugf("用户 %v 10秒内重复欢迎已被过滤", g.Uid)
			} else {
				PushToBulletSender(g.Msg)
				logx.Debug(g.Msg)
				interractGiver.interractFilter[g.Uid] = time.Now()
			}
			interractGiver.locked.Unlock()
			logx.Debugf("用户%v 已进入重复过滤列表", g.Uid)
		}
	}
END:
}
