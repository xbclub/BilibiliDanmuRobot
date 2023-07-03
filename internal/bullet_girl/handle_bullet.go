package bullet_girl

import (
	"bili_danmaku/internal/svc"
	entity "bili_danmaku/internal/types"
	"bytes"
	"compress/zlib"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"math/rand"
	"sort"
	"strings"
	"time"
)

var handler *BulletHandler

var otherSideUid map[int64]bool = make(map[int64]bool)
var topUid map[int64]bool = make(map[int64]bool)

type BulletHandler struct {
	BulletChan chan []byte
}

func pushToBulletHandler(message []byte) {
	handler.BulletChan <- message
}

func HandleBullet(ctx context.Context, svcCtx *svc.ServiceContext) {
	handler = &BulletHandler{
		BulletChan: make(chan []byte, 1000),
	}

	var message []byte
	for {
		select {
		case <-ctx.Done():
			goto END
		case message = <-handler.BulletChan:
			handle(message, svcCtx)
		}
	}
END:
}

func in(target string, src []string) bool {
	if src != nil {
		sort.Strings(src)
		index := sort.SearchStrings(src, target)
		return index < len(src) && src[index] == target
	}
	return false
}

func inWide(target string, src []string) bool {
	if src != nil {
		for _, s := range src {
			if strings.Contains(target, s) {
				return true
			}
		}
	}
	return false
}

func handle(message []byte, svcCtx *svc.ServiceContext) {
	var err error
	// 一个正文可能包含多个数据包，需要逐个解析
	index := 0
	for index < len(message) {

		// 读出包长
		var length uint32
		if err = binary.Read(bytes.NewBuffer(message[index:index+headLengthOffset]), binary.BigEndian, &length); err != nil {
			logx.Errorf("解析包长度失败", err)
			return
		}

		// 读出正文协议版本
		var ver Version
		if err = binary.Read(bytes.NewBuffer(message[index+versionOffset:index+opcodeOffset]), binary.BigEndian, &ver); err != nil {
			logx.Errorf("解析正文协议版本失败", err)
			return
		}

		// 读出操作码
		var op Opcode
		if err = binary.Read(bytes.NewBuffer(message[index+opcodeOffset:index+magicOffset]), binary.BigEndian, &op); err != nil {
			logx.Errorf("解析操作码失败", err)
			return
		}

		// 读出正文内容
		body := message[index+packageLength : index+int(length)]

		// 解析正文内容
		switch ver {
		case normalJson:
			text := &entity.CmdText{}
			_ = json.Unmarshal(body, text)
			//logx.Infof("普通json包：%s,%v,%v", text.Cmd, ver, op)
			if op == command {
				switch Cmd(text.Cmd) {

				// 处理弹幕
				case DanmuMsg:
					danmu := &entity.DanmuMsgText{}
					_ = json.Unmarshal(body, danmu)
					from := danmu.Info[2].([]interface{})

					uid := fmt.Sprintf("%.0f", from[0].(float64))
					// 如果发现弹幕在@我，那么调用机器人进行回复
					y, content := checkIsAtMe(danmu.Info[1].(string), uid, svcCtx)
					if y && danmu.Info[1].(string) != svcCtx.Config.EntryMsg {
						PushToBulletRobot(content)
					}

					logx.Infof("%.0f %v:%v", from[0].(float64), from[1], danmu.Info[1])

				// 进场特效欢迎
				case entryEffect:
					entry := &entity.EntryEffectText{}
					_ = json.Unmarshal(body, entry)
					if v, ok := svcCtx.Config.WelcomeString[fmt.Sprint(entry.Data.Uid)]; svcCtx.Config.WelcomeSwitch && ok && svcCtx.Config.EntryEffect {
						PushToBulletSender(v)
					} else if svcCtx.Config.EntryEffect {
						//PushToBulletSender(welcomeCaptain(entry.Data.CopyWriting))
						pushToInterractChan(&InterractData{
							Uid: entry.Data.Uid,
							Msg: welcomeCaptain(entry.Data.CopyWriting),
						})
					}

				// 欢迎进入房间（该功能会欢迎所有进入房间的人，可能会造成刷屏）
				case interactWord:
					interact := &entity.InteractWordText{}
					_ = json.Unmarshal(body, interact)
					// 1 进场 2 关注
					if interact.Data.MsgType == 1 {
						if v, ok := svcCtx.Config.WelcomeString[fmt.Sprint(interact.Data.Uid)]; svcCtx.Config.WelcomeSwitch && ok {
							PushToBulletSender(v)
						} else if svcCtx.Config.InteractWord {
							// 不在黑名单才欢迎
							if !inWide(interact.Data.Uname, svcCtx.Config.WelcomeBlacklistWide) &&
								!in(interact.Data.Uname, svcCtx.Config.WelcomeBlacklist) {
								pushToInterractChan(&InterractData{
									Uid: interact.Data.Uid,
									Msg: handleInterract(interact.Data.Uid, welcomeInteract(interact.Data.Uname), svcCtx),
								})
							}
						}
					} else if interact.Data.MsgType == 2 {
						msg := "感谢 " + shortName(interact.Data.Uname, 8) + " 的关注!"
						PushToBulletSender(msg)
						if svcCtx.Config.FocusDanmu != nil && len(svcCtx.Config.FocusDanmu) > 0 {
							rand.Seed(time.Now().UnixMicro())
							PushToBulletSender(svcCtx.Config.FocusDanmu[rand.Intn(len(svcCtx.Config.FocusDanmu))])
						}
					} else if interact.Data.MsgType == 3 {
						msg := "感谢 " + shortName(interact.Data.Uname, 8) + " 的分享!"
						PushToBulletSender(msg)
						if svcCtx.Config.FocusDanmu != nil && len(svcCtx.Config.FocusDanmu) > 0 {
							rand.Seed(time.Now().UnixMicro())
							PushToBulletSender(svcCtx.Config.FocusDanmu[rand.Intn(len(svcCtx.Config.FocusDanmu))])
						}
					} else {
						logx.Info(">>>>>>>>>>>>> 未识别的类型:", string(body))
					}

				// 感谢礼物
				case sendGift:
					if svcCtx.Config.ThanksGift {
						send := &entity.SendGiftText{}
						_ = json.Unmarshal(body, send)
						pushToGiftChan(send)
					}
				case "PK_BATTLE_START_NEW", "PK_BATTLE_START":
					if svcCtx.Config.PKNotice {
						info := &entity.PKStartInfo{}
						roomid := 0
						err := json.Unmarshal(body, info)
						if err != nil {
							logx.Error(err)
							logx.Errorf("pk数据解析失败:%s", string(body))
							return
						}
						if info.Data.InitInfo.RoomId == svcCtx.Config.RoomId {
							roomid = info.Data.MatchInfo.RoomId
						} else {
							roomid = info.Data.InitInfo.RoomId
						}
						logx.Debug("开始pk")
						//go handlerPK(svcCtx, body)
						if roomid == 0 {
							logx.Error("未获取的pk对手信息")
						} else {
							pushToPKChan(&roomid)
						}

					}

				case "PK_BATTLE_END", "PK_END", "PK_BATTLE_CRIT":
					// 清空串门列表
					for k := range otherSideUid {
						delete(otherSideUid, k)
					}

					//default:
					//	logx.Debug("---------------------")
					//	logx.Debug(text.Cmd)
					//	logx.Debug(string(body))
					//	logx.Debug("---------------------")
				}
			}
		case heartOrCertification:
			logx.Infof("心跳回复包")
		case normalZlib:
			b := bytes.NewReader(body)
			r, _ := zlib.NewReader(b)
			var out bytes.Buffer
			_, _ = io.Copy(&out, r)
			handle(out.Bytes(), svcCtx) // zlib解压后再进行格式解析
		}
		index += int(length)
	}
}

// 欢迎舰长语句
func welcomeCaptain(s string) string {
	s = strings.Replace(s, "\u003c%", "", 1)
	s = strings.Replace(s, "%\u003e", "", 1)

	return s
}

func welcomeInteract(name string) string {
	if strings.Contains(name, "欢迎") {
		name = strings.Replace(name, "欢迎", "", 1)
		return name
	} else {
		return name
	}
}

func shortName(uname string, alreadyLen int) string {
	s := []rune(uname)
	maxLen := (20 - alreadyLen)
	if len(s) > maxLen && maxLen > 0 {
		return string(s[0:maxLen])
	} else {
		return uname
	}
}

func handleInterract(uid int64, uname string, svcCtx *svc.ServiceContext) string {
	s := []rune(uname)
	rand.Seed(time.Now().UnixMicro())
	r := "{user}"
	if _, ook := otherSideUid[uid]; ook {
		szWelcom := "欢迎  过来串门~"
		maxLen := (20 - len([]rune(szWelcom)))
		if len(s) > maxLen && maxLen > 0 {
			return "欢迎 " + string(s[0:maxLen]) + " 过来串门~"
		} else {
			return "欢迎 " + uname + " 过来串门~"
		}
	} else {
		szWelcomOrig := svcCtx.Config.WelcomeDanmu[rand.Intn(len(svcCtx.Config.WelcomeDanmu))]
		szWelcomTmp := strings.ReplaceAll(szWelcomOrig, r, "")
		maxLen := (20 - len([]rune(szWelcomTmp)))
		logx.Info(szWelcomTmp, " ", maxLen)
		if len(s) > maxLen && maxLen > 0 {
			return strings.ReplaceAll(szWelcomOrig, r, string(s[0:maxLen]))
		} else {
			return strings.ReplaceAll(szWelcomOrig, r, uname)
		}
	}
}
