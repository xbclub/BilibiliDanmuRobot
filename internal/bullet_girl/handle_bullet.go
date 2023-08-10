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
			logx.Errorf("解析包长度失败:%v", err)
			return
		}

		// 读出正文协议版本
		var ver Version
		if err = binary.Read(bytes.NewBuffer(message[index+versionOffset:index+opcodeOffset]), binary.BigEndian, &ver); err != nil {
			logx.Errorf("解析正文协议版本失败:%v", err)
			return
		}

		// 读出操作码
		var op Opcode
		if err = binary.Read(bytes.NewBuffer(message[index+opcodeOffset:index+magicOffset]), binary.BigEndian, &op); err != nil {
			logx.Errorf("解析操作码失败:%v", err)
			return
		}

		// 读出正文内容
		body := message[index+packageLength : index+int(length)]

		// 解析正文内容
		switch ver {
		case normalJson:
			logx.Debug(string(body))
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

					// @帮助 打出来关键词
					if strings.Compare("@帮助", danmu.Info[1].(string)) == 0 {
						s := fmt.Sprintf("发送带有 %s 的弹幕和我互动", svcCtx.Config.TalkRobotCmd)
						logx.Info(s)
						PushToBulletSender(" ")
						PushToBulletSender(s)
						PushToBulletSender("请尽情调戏我吧!")
					}

					uid := fmt.Sprintf("%.0f", from[0].(float64))
					// 如果发现弹幕在@我，那么调用机器人进行回复
					y, content := checkIsAtMe(danmu.Info[1].(string), uid, svcCtx)
					if y && len(content) > 0 && danmu.Info[1].(string) != svcCtx.Config.EntryMsg {
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
								if svcCtx.Config.InteractWordByTime {
									msg := handleInterractByTime(interact.Data.Uid, welcomeInteract(interact.Data.Uname), svcCtx)
									ms := strings.Split(msg, "\n")
									if len(ms) > 1 {
										for i, s := range ms {
											pushToInterractChan(&InterractData{
												Uid: interact.Data.Uid + int64(i),
												Msg: s,
											})
										}
									} else {
										pushToInterractChan(&InterractData{
											Uid: interact.Data.Uid,
											Msg: msg,
										})
									}
								} else {
									msg := handleInterract(interact.Data.Uid, welcomeInteract(interact.Data.Uname), svcCtx)
									ms := strings.Split(msg, "\n")
									if len(ms) > 1 {
										for i, s := range ms {
											pushToInterractChan(&InterractData{
												Uid: interact.Data.Uid + int64(i),
												Msg: s,
											})
										}
									} else {
										pushToInterractChan(&InterractData{
											Uid: interact.Data.Uid,
											Msg: msg,
										})
									}
								}
							}
						}
					} else if interact.Data.MsgType == 2 {
						if svcCtx.Config.InteractWord {
							msg := "感谢 " + shortName(interact.Data.Uname, 8, svcCtx.Config.DanmuLen) + " 的关注!"
							PushToBulletSender(msg)
							if svcCtx.Config.FocusDanmu != nil && len(svcCtx.Config.FocusDanmu) > 0 {
								rand.Seed(time.Now().UnixMicro())
								PushToBulletSender(svcCtx.Config.FocusDanmu[rand.Intn(len(svcCtx.Config.FocusDanmu))])
							}
						}
					} else if interact.Data.MsgType == 3 {
						if svcCtx.Config.InteractWord {
							msg := "感谢 " + shortName(interact.Data.Uname, 8, svcCtx.Config.DanmuLen) + " 的分享!"
							PushToBulletSender(msg)
							if svcCtx.Config.FocusDanmu != nil && len(svcCtx.Config.FocusDanmu) > 0 {
								rand.Seed(time.Now().UnixMicro())
								PushToBulletSender(svcCtx.Config.FocusDanmu[rand.Intn(len(svcCtx.Config.FocusDanmu))])
							}
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

func shortName(uname string, alreadyLen, danmuLen int) string {
	s := []rune(uname)
	maxLen := (danmuLen - alreadyLen)
	if len(s) > maxLen && maxLen > 0 {
		return string(s[0:maxLen-1]) + "…"
	} else {
		return uname
	}
}

func handleInterractByTime(uid int64, uname string, svcCtx *svc.ServiceContext) string {

	if _, ook := otherSideUid[uid]; ook {
		return handleInterract(uid, uname, svcCtx)
	}
	// 凌晨 - Early morning   2:00--5:00
	// 早晨 - Morning   5:00--9:00
	// 上午 - Late morning / Mid-morning  9:00--11:00
	// 中午 - Noon  11:00--14:00
	// 下午 - Afternoon 14:00 -- 20:00
	// 晚上 - Evening / Night 20:00--00:00
	// 午夜 - Midnight 00:00 -- 2:00
	//s := []rune(uname)
	rand.Seed(time.Now().UnixMicro())
	r := "{user}"

	if svcCtx.Config.InteractWordByTime &&
		svcCtx.Config.WelcomeDanmuByTime != nil &&
		len(svcCtx.Config.WelcomeDanmuByTime) > 0 {

		now := time.Now().Hour()

		key := ""
		switch now {
		case 0, 1:
			// 午夜
			key = "midnight"

		case 2, 3, 4:
			// 凌晨
			key = "earlymorning"

		case 5, 6, 7, 8:
			// 早上
			key = "morning"

		case 9, 10:
			// 上午
			key = "latemorning"

		case 11, 12, 13:
			// 中午
			key = "noon"

		case 14, 15, 16, 17, 18, 19:
			// 下午
			key = "afternoon"

		case 20, 21, 22, 23:
			// 晚上
			key = "night"
		}

		for _, danmuCfg := range svcCtx.Config.WelcomeDanmuByTime {
			if danmuCfg.Key == key {
				if danmuCfg.Enabled && len(danmuCfg.Danmu) > 0 {
					szWelcomOrig := danmuCfg.Danmu[rand.Intn(len(danmuCfg.Danmu))]

					welcome := strings.ReplaceAll(szWelcomOrig, r, shortName(uname, 3, svcCtx.Config.DanmuLen))
					rWelcome := []rune(welcome)
					if len(rWelcome) > svcCtx.Config.DanmuLen {
						szWelcomTmp := strings.ReplaceAll(szWelcomOrig, r+", ", r+"\n")
						szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+",", r+"\n")
						szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+"，", r+"\n")
						return strings.ReplaceAll(szWelcomTmp, r, uname)
					} else {
						return welcome
					}
				} else {
					return handleInterract(uid, uname, svcCtx)
				}
			}
		}
		return handleInterract(uid, uname, svcCtx)
	} else {
		return handleInterract(uid, uname, svcCtx)
	}
}

func handleInterract(uid int64, uname string, svcCtx *svc.ServiceContext) string {
	s := []rune(uname)
	rand.Seed(time.Now().UnixMicro())
	r := "{user}"
	if _, ook := otherSideUid[uid]; ook {
		szWelcom := "欢迎  过来串门~"
		maxLen := (svcCtx.Config.DanmuLen - len([]rune(szWelcom)))
		if len(s) > maxLen && maxLen > 0 {
			return "欢迎 " + string(s[0:maxLen-1]) + "… 过来串门~"
		} else {
			return "欢迎 " + uname + " 过来串门~"
		}
	} else {
		szWelcomOrig := svcCtx.Config.WelcomeDanmu[rand.Intn(len(svcCtx.Config.WelcomeDanmu))]

		welcome := strings.ReplaceAll(szWelcomOrig, r, shortName(uname, 3, svcCtx.Config.DanmuLen))
		rWelcome := []rune(welcome)
		if len(rWelcome) > svcCtx.Config.DanmuLen {
			szWelcomTmp := strings.ReplaceAll(szWelcomOrig, r+", ", r+"\n")
			szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+",", r+"\n")
			szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+"，", r+"\n")
			return strings.ReplaceAll(szWelcomTmp, r, uname)
		} else {
			return welcome
		}
	}
}
