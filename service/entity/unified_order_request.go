package entity

type UnifiedOrderRequest struct {
	Body     string `json:"body"`
	TotalFee uint32 `json:"totalFee"`
	UserId   string `json:"userId"`
}
