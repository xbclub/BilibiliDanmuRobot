package bullet_girl

import (
	"bili_danmaku/internal/http"
	"bili_danmaku/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

var robot *BulletRobot

type BulletRobot struct {
	bulletRobotChan chan string
}

func PushToBulletRobot(content string) {
	logx.Infof("PushToBulletRobot成功：%s", content)
	robot.bulletRobotChan <- content
}

func StartBulletRobot(ctx context.Context, svcCtx *svc.ServiceContext) {
	robot = &BulletRobot{
		bulletRobotChan: make(chan string, 1000),
	}

	var content string

	for {
		select {
		case <-ctx.Done():
			goto END
		case content = <-robot.bulletRobotChan:
			handleRobotBullet(content, svcCtx)
		}
	}
END:
}

func handleRobotBullet(content string, svcCtx *svc.ServiceContext) {
	var err error
	var reply string
	if svcCtx.Config.RobotMode == "ChatGPT" {
		if reply, err = http.RequestChatgptRobot(content, svcCtx); err != nil {
			logx.Errorf("请求机器人失败：", err)
			PushToBulletSender("不好意思，机器人坏掉了...")
			return
		}
	} else {
		if reply, err = http.RequestQingyunkeRobot(content); err != nil {
			logx.Errorf("请求机器人失败：", err)
			PushToBulletSender("不好意思，机器人坏掉了...")
			return
		}
		bulltes := splitRobotReply(reply, svcCtx)
		for _, v := range bulltes {
			PushToBulletSender(v)
		}
		return
	}
	PushToBulletSender(reply)
	logx.Infof("机器人回复：%s", reply)

}

// 将机器人回复语句中的 {br} 进行分割
// b站弹幕一次只能发20个字符，需要切分
func splitRobotReply(content string, svcCtx *svc.ServiceContext) []string {

	// 将机器人回复中的菲菲替换为橘子
	content = strings.ReplaceAll(content, "菲菲", svcCtx.Config.RobotName)

	//var res []string
	reply := strings.Split(content, "{br}")

	//for _, r := range reply {
	//	// 长度大于20再分割
	//	zh := []rune(r)
	//	if len(zh) > 20 {
	//		i := 0
	//		for i < len(zh) {
	//			if i+20 > len(zh) {
	//				res = append(res, string(zh[i:]))
	//			} else {
	//				res = append(res, string(zh[i:i+20]))
	//			}
	//			i += 20
	//		}
	//	} else {
	//		res = append(res, string(zh))
	//	}
	//}
	return reply
}

// 检查弹幕是否在@我，返回bool和@我要说的内容
func checkIsAtMe(msg, u string, svcCtx *svc.ServiceContext) (bool, string) {
	// 自己发的包含关键字 不与理会 避免递归

	userId, ok := http.CookieList["DedeUserID"]

	if ok && userId != u && strings.Contains(msg, svcCtx.Config.TalkRobotCmd) {
		return true, strings.ReplaceAll(msg, svcCtx.Config.TalkRobotCmd, "")
	} else {
		return false, ""
	}
}
