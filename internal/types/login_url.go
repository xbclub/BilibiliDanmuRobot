package types

type LoginUrl struct {
	Data struct {
		Url      string `json:"url"`
		OauthKey string `json:"oauthKey"`
	} `json:"data"`
}

type LoginInfoPre struct {
	Status bool `json:"status"`
}

type LoginInfoData struct {
	Data struct {
		Url          string `json:"url"`
		RefreshToken string `json:"refresh_token"`
	} `json:"data"`
}

type LoginInfoCookies struct {
	SetCookie []string `json:"Set-Cookie"`
}
