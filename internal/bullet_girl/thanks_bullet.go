package bullet_girl

import (
	entity "bili_danmaku/internal/types"
	"context"
	"fmt"
	"sync"
	"time"
)

// 检测到礼物，push [uname]->[giftName]->[cost]，number+1
// 每10s统计一次礼物，并进行感谢，礼物价值高于x元加一句大气

var thanksGiver *GiftThanksGiver

type GiftThanksGiver struct {
	giftTable map[string]map[string]map[string]int
	tableMu   sync.RWMutex
	giftChan  chan *entity.SendGiftText
}

func pushToGiftChan(g *entity.SendGiftText) {
	thanksGiver.giftChan <- g
}

func ThanksGift(ctx context.Context) {

	thanksGiver = &GiftThanksGiver{
		giftTable: make(map[string]map[string]map[string]int),
		tableMu:   sync.RWMutex{},
		giftChan:  make(chan *entity.SendGiftText, 1000),
	}

	var g *entity.SendGiftText
	var w = 10 * time.Second
	var t = time.NewTimer(w)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			goto END
		case <-t.C:
			summarizeGift()
			t.Reset(w)
		case g = <-thanksGiver.giftChan:
			if _, ok := thanksGiver.giftTable[g.Data.Uname]; !ok {
				thanksGiver.giftTable[g.Data.Uname] = make(map[string]map[string]int)
			}
			if _, ok := thanksGiver.giftTable[g.Data.Uname][g.Data.GiftName]; !ok {
				thanksGiver.giftTable[g.Data.Uname][g.Data.GiftName] = make(map[string]int)
			}
			thanksGiver.giftTable[g.Data.Uname][g.Data.GiftName]["cost"] += g.Data.Price
			thanksGiver.giftTable[g.Data.Uname][g.Data.GiftName]["count"] += 1
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
		for k, v := range giftstring {
			if k == 0 {
				msg += "感谢 " + name + " 的 " + v
			} else {
				msg += "，" + v
			}
		}
		PushToBulletSender(msg)
		//fmt.Println("礼物-----", name, giftstring)
		// 总打赏高于x元，加一句大气
		if sumCost >= 50000 { // 2元
			PushToBulletSender(name + "老板大气大气")
		}
		delete(thanksGiver.giftTable, name)
	}
}
