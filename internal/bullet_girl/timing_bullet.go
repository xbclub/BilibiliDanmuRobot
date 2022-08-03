package bullet_girl

import (
	"context"
	"github.com/gorhill/cronexpr"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

const (
	Save = iota
	Remove
)

// 调度表
var scheduler *BulletTaskScheduler
var id = 0

// 弹幕
type Bullet struct {
	id   int
	msg  string
	expr string // 定时crontab表达式
}

// 弹幕定时任务
type BulletTask struct {
	bullet *Bullet
	expr   *cronexpr.Expression
	next   time.Time // 下次调度时间
}

// 弹幕事件，删除/创建 定时弹幕
type BulletEvent struct {
	spec       int
	bulletTask *BulletTask
}

// 定时弹幕任务调度表
type BulletTaskScheduler struct {
	table     map[int]*BulletTask
	eventChan chan *BulletEvent
}

func NewBullet(msg string, expr string) *Bullet {
	id++
	return &Bullet{id: id, msg: msg, expr: expr}
}

func NewBulletTask(b *Bullet) *BulletTask {
	bt := &BulletTask{}
	bt.bullet = b
	bt.expr = cronexpr.MustParse(b.expr)
	bt.next = bt.expr.Next(time.Now())
	return bt
}

func NewBulletEvent(spec int, bt *BulletTask) *BulletEvent {
	return &BulletEvent{
		spec:       spec,
		bulletTask: bt,
	}
}

func PushToBulletEvent(be *BulletEvent) {
	logx.Infof("PushBulletEvent成功", be.bulletTask.bullet.msg)
	scheduler.eventChan <- be
}

// 定时弹幕任务调度
func StartTimingBullet(ctx context.Context) {

	// 初始化任务表
	scheduler = &BulletTaskScheduler{
		table:     make(map[int]*BulletTask),
		eventChan: make(chan *BulletEvent, 1000),
	}

	var be *BulletEvent

	interval := CalculateAndRun()
	t := time.NewTimer(interval)

	defer t.Stop()

	for {
		select {
		// 事件处理
		case be = <-scheduler.eventChan:
			HandleBulletEvent(be)
		// 关闭goroutine
		case <-ctx.Done():
			goto END
		// 到达等待时间，开始执行定时任务
		case <-t.C:
		}
		interval = CalculateAndRun()
		t.Reset(interval)
	}

END:
}

// 定时弹幕事件处理
func HandleBulletEvent(be *BulletEvent) {
	switch be.spec {
	case Save:
		logx.Infof("task保存成功", be.bulletTask.bullet.msg)
		scheduler.table[be.bulletTask.bullet.id] = be.bulletTask
	case Remove:
		delete(scheduler.table, be.bulletTask.bullet.id)
	}
}

// 在所有定时任务中计算出需要等待的时间，并执行到期任务
func CalculateAndRun() time.Duration {

	var interval *time.Time
	now := time.Now()

	for _, bt := range scheduler.table {

		// 执行到期任务
		if bt.next.Before(now) || bt.next.Equal(now) {
			PushToBulletSender(bt.bullet.msg)
			bt.next = bt.expr.Next(now) // 更新下一次执行时间
		}

		// 确定最近任务间隔时间
		if interval == nil || bt.next.Before(*interval) {
			interval = &bt.next
		}
	}

	// 没有任务固定等待1s
	if interval == nil {
		return time.Second
	}

	return (*interval).Sub(time.Now())
}
