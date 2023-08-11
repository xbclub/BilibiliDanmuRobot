package bullet_girl

import (
	"bili_danmaku/internal/http"
	"bili_danmaku/internal/svc"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	http2 "net/http"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type Opcode uint32  // 数据包业务类型
type Version uint16 // 正文类型及压缩方式
type Cmd string     // 命令类型

const (
	normalJson           Version = 0 // 正文为json格式的弹幕
	heartOrCertification Version = 1 // 心跳或认证包正文不压缩，客户端发送的心跳包无正文，服务队发送的心跳包正文为4字节数据，表示人气值
	normalZlib           Version = 2 // 普通包正文使用zlib压缩
	//normalBrotli         Version = 3               // 普通包正文使用brotli压缩，解压后为一个普通包（头部协议为0），需要再次解析出正文
	heartBeat     Opcode = 2 // 心跳包
	command       Opcode = 5 // 命令包
	certification Opcode = 7 // 认证包
	//enterRoom            Opcode  = 8               // 进入房间
	DanmuMsg Cmd = "DANMU_MSG" // 弹幕消息
	//welcomeGuard         Cmd     = "WELCOME_GUARD" // 欢迎xxx老爷
	entryEffect Cmd = "ENTRY_EFFECT" // 欢迎舰长进入房间
	//welcome              Cmd     = "WELCOME"       // 欢迎xxx进入房间
	interactWord Cmd = "INTERACT_WORD" // 进入房间
	sendGift     Cmd = "SEND_GIFT"     // 发现送礼物
)

// 关于数据包格式的常量
const (
	packageLength = 16 // 包长度
	magicNumber   = 1  // 包头最后的魔数

	// 包头中，字节位置偏移量
	headLengthOffset = 4
	versionOffset    = 6
	opcodeOffset     = 8
	magicOffset      = 12
)

type CertificationPackageBody struct {
	RoomId   int    `json:"roomid"`
	Uid      int64  `json:"uid"`
	Key      string `json:"key"`
	Protover int8   `json:"protover"`
	Buvid    string `json:"buvid,optional"`
	Platform string `json:"platform,default=web"`
	Type     int8   `json:"type,default=2"`
}

// 生成数据包头部
func GeneratePackageHead(bodyLength uint32, opcode Opcode) ([]byte, error) {
	var err error
	head := bytes.NewBuffer([]byte{})

	// 总长度 该值占4字节
	if err = binary.Write(head, binary.BigEndian, bodyLength+uint32(packageLength)); err != nil {
		return nil, err
	}
	// 头部长度 固定16 该值占2字节
	if err = binary.Write(head, binary.BigEndian, uint16(packageLength)); err != nil {
		return nil, err
	}
	// 协议版本号 固定1 该值占2字节
	if err = binary.Write(head, binary.BigEndian, heartOrCertification); err != nil {
		return nil, err
	}
	// 操作码 该值占4字节
	if err = binary.Write(head, binary.BigEndian, opcode); err != nil {
		return nil, err
	}
	// 包序号 可取常数1 该值占4字节
	if err = binary.Write(head, binary.BigEndian, uint32(magicNumber)); err != nil {
		return nil, err
	}

	return head.Bytes(), nil
}

// 生成请求数据包，由包头和正文组成
func GenerateCertificationPackage(svcCtx *svc.ServiceContext, token string, uid int64) ([]byte, error) {
	var err error
	var head []byte
	var body []byte

	cpb := &CertificationPackageBody{
		RoomId: svcCtx.Config.RoomId,
		Key:    token,
		Uid:    uid,
		//Protover: 3,
	}
	body, _ = json.Marshal(cpb)
	logx.Debug(string(body))
	if head, err = GeneratePackageHead(uint32(len(body)), certification); err != nil {
		logx.Errorf("生成包头失败：%v", err)
	}

	return append(head[:], body[:]...), nil
}

// 30s发送一次心跳包
func StartHeartBeat(ctx context.Context, conn *websocket.Conn) {
	var hb []byte
	var err error
	t := time.NewTimer(30 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			goto END
		case <-t.C:
			t.Reset(30 * time.Second)
			// 心跳包无正文
			if hb, err = GeneratePackageHead(0, heartBeat); err != nil {
				logx.Errorf("心跳包组装错误：%v", err)
			}
			if err = conn.WriteMessage(websocket.BinaryMessage, hb); err != nil {
				logx.Errorf("发送心跳包失败：%v", err)
				return
			}
		}
	}
END:
}

func StartCatchBullet(ctx context.Context, svcCtx *svc.ServiceContext) {
	var err error
	var cert []byte
	var conn *websocket.Conn
	var message []byte
	token, err := http.GetDanmuToken(svcCtx.Config.RoomId)
	if err != nil {
		logx.Errorf("获取直播信息流认证失败：%v", err)
		return
	}
	user := http.GetUserInfo()
	if user == nil {
		logx.Error("获取登录用户uid失败")
		return
	}
	// 连接ws服务器
	//logx.Debug(fmt.Sprintf("wss://%v:%v/sub", token.Data.HostList[0].Host, token.Data.HostList[0].WssPort))
	//if conn, _, err = websocket.DefaultDialer.Dial(fmt.Sprintf("wss://%v:%v/sub", token.Data.HostList[0].Host, token.Data.HostList[0].WssPort), nil); err != nil {
	header := http2.Header{}
	header.Set("Origin", "https://live.bilibili.com")
	header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.200")
	if conn, _, err = websocket.DefaultDialer.Dial(svcCtx.Config.WsServerUrl, header); err != nil {
		logx.Errorf("websocket连接失败：%v", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			logx.Error(err)
		}
	}(conn)

	// 组装认证包
	if cert, err = GenerateCertificationPackage(svcCtx, token.Data.Token, user.Uid); err != nil {
		logx.Errorf("组装认证包错误：%v", err)
		return
	}

	// 发送认证包
	if err = conn.WriteMessage(websocket.BinaryMessage, cert); err != nil {
		logx.Errorf("发送认证包失败：%v", err)
		return
	}

	// 开启心跳包
	hbCtx, hbCancel := context.WithCancel(context.Background())
	defer hbCancel()
	go StartHeartBeat(hbCtx, conn)

	// 循环接受信息
	for {
		select {
		case <-ctx.Done():
			hbCancel()
			goto END
		default:
			if _, message, err = conn.ReadMessage(); err != nil {
				logx.Errorf("websocket读取消息失败：%v", err)
				continue
			}
			pushToBulletHandler(message)
		}
	}
END:
}
