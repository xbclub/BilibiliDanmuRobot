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
		Uid        int64 `json:"uid"`
		LiveStatus int   `json:"live_status"`
	} `json:"data"`
}

type Userinfo struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    struct {
		Info struct {
			Uid            int64  `json:"uid"`
			Uname          string `json:"uname"`
			Face           string `json:"face"`
			OfficialVerify struct {
				Type int    `json:"type"`
				Desc string `json:"desc"`
			} `json:"official_verify"`
			Gender int `json:"gender"`
		} `json:"info"`
		Exp struct {
			MasterLevel struct {
				Level   int   `json:"level"`
				Color   int   `json:"color"`
				Current []int `json:"current"`
				Next    []int `json:"next"`
			} `json:"master_level"`
		} `json:"exp"`
		FollowerNum  int    `json:"follower_num"`
		RoomId       int    `json:"room_id"`
		MedalName    string `json:"medal_name"`
		GloryCount   int    `json:"glory_count"`
		Pendant      string `json:"pendant"`
		LinkGroupNum int    `json:"link_group_num"`
		RoomNews     struct {
			Content   string `json:"content"`
			Ctime     string `json:"ctime"`
			CtimeText string `json:"ctime_text"`
		} `json:"room_news"`
	} `json:"data"`
}

type TopListInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Info struct {
			Num                     int `json:"num"`
			Page                    int `json:"page"`
			Now                     int `json:"now"`
			AchievementLevel        int `json:"achievement_level"`
			AnchorGuardAchieveLevel int `json:"anchor_guard_achieve_level"`
		} `json:"info"`
		List []struct {
			Uid           int64  `json:"uid"`
			Ruid          int64  `json:"ruid"`
			Rank          int    `json:"rank"`
			Username      string `json:"username"`
			Face          string `json:"face"`
			IsAlive       int    `json:"is_alive"`
			GuardLevel    int    `json:"guard_level"`
			GuardSubLevel int    `json:"guard_sub_level"`
			MedalInfo     struct {
				MedalName        string `json:"medal_name"`
				MedalLevel       int    `json:"medal_level"`
				MedalColorStart  int    `json:"medal_color_start"`
				MedalColorEnd    int    `json:"medal_color_end"`
				MedalColorBorder int    `json:"medal_color_border"`
			} `json:"medal_info"`
		} `json:"list"`
		Top3 []struct {
			Uid           int64  `json:"uid"`
			Ruid          int64  `json:"ruid"`
			Rank          int    `json:"rank"`
			Username      string `json:"username"`
			Face          string `json:"face"`
			IsAlive       int    `json:"is_alive"`
			GuardLevel    int    `json:"guard_level"`
			GuardSubLevel int    `json:"guard_sub_level"`
			MedalInfo     struct {
				MedalName        string `json:"medal_name"`
				MedalLevel       int    `json:"medal_level"`
				MedalColorStart  int    `json:"medal_color_start"`
				MedalColorEnd    int    `json:"medal_color_end"`
				MedalColorBorder int    `json:"medal_color_border"`
			} `json:"medal_info"`
		} `json:"top3"`
		MyFollowInfo struct {
			GuardLevel    int    `json:"guard_level"`
			AccompanyDays int    `json:"accompany_days"`
			ExpiredTime   string `json:"expired_time"`
			AutoRenew     int    `json:"auto_renew"`
			RenewRemind   struct {
				Content string `json:"content"`
				Type    int    `json:"type"`
				Hint    string `json:"hint"`
			} `json:"renew_remind"`
			MedalInfo struct {
				MedalName        string `json:"medal_name"`
				MedalLevel       int    `json:"medal_level"`
				MedalColorStart  int    `json:"medal_color_start"`
				MedalColorEnd    int    `json:"medal_color_end"`
				MedalColorBorder int    `json:"medal_color_border"`
			} `json:"medal_info"`
			Rank int    `json:"rank"`
			Ruid int    `json:"ruid"`
			Face string `json:"face"`
		} `json:"my_follow_info"`
		GuardWarn struct {
			IsWarn      int    `json:"is_warn"`
			Warn        string `json:"warn"`
			Expired     int    `json:"expired"`
			WillExpired int    `json:"will_expired"`
			Address     string `json:"address"`
		} `json:"guard_warn"`
		ExistBenefit  bool   `json:"exist_benefit"`
		RemindBenefit string `json:"remind_benefit"`
	} `json:"data"`
}
