package types

type UserinfoLite struct {
	Islogin  bool
	Username string
	Avactor  string
	Uid      int64
}
type UserInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		IsLogin bool   `json:"isLogin"`
		Mid     int64  `json:"mid"`
		Face    string `json:"face"`
		Uname   string `json:"uname"`
	} `json:"data"`
}
