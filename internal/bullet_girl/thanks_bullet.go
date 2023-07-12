package bullet_girl

import (
	"bili_danmaku/internal/svc"
	entity "bili_danmaku/internal/types"
	"context"
	"fmt"
	"sync"
	"time"
)

// 检测到礼物，push [uname]->[giftName]->[cost]，number+1
// 每3s统计一次礼物，并进行感谢，礼物价值高于x元加一句大气

var thanksGiver *GiftThanksGiver

type GiftThanksGiver struct {
	giftTable map[string]map[string]map[string]int
	locked    *sync.Mutex
	tableMu   sync.RWMutex
	giftChan  chan *entity.SendGiftText
}

func pushToGiftChan(g *entity.SendGiftText) {
	thanksGiver.giftChan <- g
}

func ThanksGift(ctx context.Context, svcCtx *svc.ServiceContext) {

	thanksGiver = &GiftThanksGiver{
		giftTable: make(map[string]map[string]map[string]int),
		locked:    new(sync.Mutex),
		tableMu:   sync.RWMutex{},
		giftChan:  make(chan *entity.SendGiftText, 1000),
	}

	var g *entity.SendGiftText
	var w = time.Duration(svcCtx.Config.ThanksGiftTimeout) * time.Second
	var t = time.NewTimer(w)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			goto END
		case <-t.C:
			thanksGiver.locked.Lock()
			summarizeGift()
			thanksGiver.locked.Unlock()
			t.Reset(w)
		case g = <-thanksGiver.giftChan:
			thanksGiver.locked.Lock()
			if _, ok := thanksGiver.giftTable[g.Data.Uname]; !ok {
				thanksGiver.giftTable[g.Data.Uname] = make(map[string]map[string]int)
			}
			if _, ok := thanksGiver.giftTable[g.Data.Uname][g.Data.GiftName]; !ok {
				thanksGiver.giftTable[g.Data.Uname][g.Data.GiftName] = make(map[string]int)
			}
			thanksGiver.giftTable[g.Data.Uname][g.Data.GiftName]["cost"] += g.Data.Price
			thanksGiver.giftTable[g.Data.Uname][g.Data.GiftName]["count"] += g.Data.Num
			thanksGiver.locked.Unlock()
		}
	}
END:
}

func summarizeGift() {
	for name, m := range thanksGiver.giftTable {
		sumCost := 0
		giftstring := []string{}
		msg := ""
		for gift, cost := range m {
			giftstring = append(giftstring, fmt.Sprintf("%d个%s", cost["count"], gift))
			// 计算打赏金额
			sumCost += cost["cost"]

			// 感谢完后立刻清空map
			delete(m, gift)
		}
		msg = "感谢 " + name + " 的"
		PushToBulletSender(msg)
		for k, v := range giftstring {
			if k == 0 {
				msg = v
			} else {
				msg += "，" + v
			}
		}
		PushToBulletSender(msg)
		//fmt.Println("礼物-----", name, giftstring)
		// 总打赏高于x元，加一句大气
		if sumCost >= 50000 { // 50元
			PushToBulletSender(name + "老板大气大气")
		}
		delete(thanksGiver.giftTable, name)
	}
}
