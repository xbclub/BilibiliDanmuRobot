package types

// 弹幕服务器信息
type ResponseBulletInfo struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Ttl     int                 `json:"ttl"`
	Data    *ResponseBulletData `json:"data"`
}

type ResponseBulletData struct {
	BusinessId       int                   `json:"business_id"`
	Group            string                `json:"group"`
	HostList         []*ResponseBulletHost `json:"host_list"`
	MaxDelay         int                   `json:"max_delay"`
	RefreshRate      int                   `json:"refresh_rate"`
	RefreshRowFactor float64               `json:"refresh_row_factor"`
	Token            string                `json:"token"`
}

type ResponseBulletHost struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	WssPort int    `json:"wss_port"`
	WsPort  int    `json:"ws_port"`
}
