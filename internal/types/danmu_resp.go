package types

type DanmuResp struct {
	Code int `json:"code"`
	Data struct {
		ModeInfo struct {
			Mode           int    `json:"mode"`
			ShowPlayerType int    `json:"show_player_type"`
			Extra          string `json:"extra"`
		} `json:"mode_info"`
	} `json:"data"`
	Message string `json:"message"`
	Msg     string `json:"msg"`
}
