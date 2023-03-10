package entity

type GenerateUserRequest struct {
	OpenId           string `json:"openId"`
	AvatarUrl        string `json:"avatarUrl"`
	City             string `json:"city"`
	Country          string `json:"country"`
	Gender           int    `json:"gender"`
	Language         string `json:"language"`
	NickName         string `json:"nickName"`
	Province         string `json:"province"`
	AppId            string `json:"appId"`
	RemainUsageCount int64  `json:"remainUsageCount"`
}
