package bullet_girl

import (
	"bili_danmaku/internal/http"
	entity "bili_danmaku/internal/types"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

var pkGiver *PKGiver

type PKGiver struct {
	pkFilter map[int]time.Time
	locked   *sync.Mutex
	tableMu  sync.RWMutex
	pkChan   chan *int
}

func pushToPKChan(g *int) {
	pkGiver.pkChan <- g
}

func PK(ctx context.Context) {

	pkGiver = &PKGiver{
		pkFilter: map[int]time.Time{},
		locked:   new(sync.Mutex),
		tableMu:  sync.RWMutex{},
		pkChan:   make(chan *int, 1000),
	}

	var g *int
	var w = 10 * time.Second
	var t = time.NewTimer(w)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			goto END
		case <-t.C:
			if len(pkGiver.pkFilter) > 0 {
				pkGiver.locked.Lock()
				for k, v := range pkGiver.pkFilter {
					if v.Add(w).Unix() < time.Now().Unix() {
						delete(pkGiver.pkFilter, k)
						logx.Debugf("pk room %v 已从重复过滤列表移除", k)
					}
				}
				pkGiver.locked.Unlock()
			}

			t.Reset(w)
		case g = <-pkGiver.pkChan:
			//pkGiver.tmpmsg = append(pkGiver.tmpmsg, *g)
			pkGiver.locked.Lock()
			if value, ok := pkGiver.pkFilter[*g]; ok && value.Add(w).Unix() >= time.Now().Unix() {
				logx.Debugf("pk room %v 10秒内重复获取数据已被过滤", *g)
			} else {
				logx.Debugf("正在处理pk信息")
				handlerPK(*g)
				logx.Debug(*g)
				pkGiver.pkFilter[*g] = time.Now()
			}
			pkGiver.locked.Unlock()
			logx.Debugf("pk room %v 已进入重复过滤列表", *g)
		}
	}
END:
}
func handlerPK(roomid int) {

	toplist := &entity.TopListInfo{}
	toplistalive := 0
	rankcount := 0

	userinfo, err := http.Userinfo(roomid)
	if err != nil {
		logx.Error(err)
		return
	}
	toppage := 1
	listInfo, err := http.TopListInfo(roomid, userinfo.Data.Info.Uid, toppage)
	if err != nil {
		logx.Error(err)
		return
	}

	tmpPage := listInfo.Data.Info.Page
	for toppage += 1; toppage <= tmpPage; toppage++ {
		toplist, err = http.TopListInfo(roomid, userinfo.Data.Info.Uid, toppage)
		if err != nil {
			logx.Error(err)
			continue
		}
		listInfo.Data.List = append(listInfo.Data.List, toplist.Data.List...)
	}

	for _, data := range listInfo.Data.List {
		if data.IsAlive == 1 {
			toplistalive++
		}
	}
	rankListInfo, err := http.RankListInfo(roomid, userinfo.Data.Info.Uid, 1)
	if err != nil {
		logx.Error(err)
		return
	}
	for _, data := range rankListInfo.Data.OnlineRankItem {
		rankcount += data.Score
	}
	//PushToBulletSender(fmt.Sprintf("当前对手:%v，%v船，%v粉,对面有%v名船长在线，高能榜%v人，榜前50贡献%v分", userinfo.Data.Info.Uname, listInfo.Data.Info.Num, userinfo.Data.FollowerNum, toplistalive, rankListInfo.Data.OnlineNum, rankcount))
	PushToBulletSender(fmt.Sprintf("当前对手:%v", userinfo.Data.Info.Uname))
	PushToBulletSender(fmt.Sprintf("共%v船，%v粉", listInfo.Data.Info.Num, userinfo.Data.FollowerNum))
	//PushToBulletSender(fmt.Sprintf("对面有%v船在线，高能榜%v人", toplistalive, rankListInfo.Data.OnlineNum))
	PushToBulletSender(fmt.Sprintf("高能榜%v人", rankListInfo.Data.OnlineNum))
	PushToBulletSender(fmt.Sprintf("榜前50贡献%v分", rankcount))
}
