package entity

type LoginRequest struct {
	Info   string `json:"info"`
	OpenId string `json:"openid"`
}
