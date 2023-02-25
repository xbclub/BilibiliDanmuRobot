package types

type CmdText struct {
	Cmd string `json:"cmd"`
}

type DanmuMsgText struct {
	Info []interface{} `json:"info"`
}

type EntryEffectText struct {
	Data struct {
		Uid         int64  `json:"uid"`
		CopyWriting string `json:"copy_writing"`
	} `json:"data"`
}

type InteractWordText struct {
	Data struct {
		Uname   string `json:"uname"`
		Uid     int64  `json:"uid"`
		MsgType int32  `json:"msg_type"`
	} `json:"data"`
}

type SendGiftText struct {
	Data struct {
		Action   string `json:"action"`
		GiftName string `json:"giftName"`
		Uname    string `json:"uname"`
		Price    int    `json:"price"`
		Num      int    `json:"num"`
	} `json:"data"`
}
