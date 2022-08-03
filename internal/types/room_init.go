package types

const (
	NotStarted = 0 // 未开播
	Live       = 1 // 直播中
	Carousel   = 2 // 轮播中
)

type RoomInitStatus struct {
	Code int `json:"code"`
}

type RoomInitInfo struct {
	Data struct {
		LiveStatus int `json:"live_status"`
	} `json:"data"`
}
