package entity

type OpenData struct {
	Data    UserInfo `json:"data"`
	CloudID string   `json:"cloudID"`
}

type UserInfo struct {
	OpenId    string    `json:"openId"`
	AvatarUrl string    `json:"avatarUrl"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	Gender    int       `json:"gender"`
	Language  string    `json:"language"`
	NickName  string    `json:"nickName"`
	Province  string    `json:"province"`
	WaterMark WaterMark `json:"watermark"`
}

type WaterMark struct {
	Timestamp int32  `json:"timestamp"`
	AppId     string `json:"appid"`
}
