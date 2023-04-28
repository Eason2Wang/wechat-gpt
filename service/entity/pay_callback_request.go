package entity

type PayCallbackRequest struct {
	ReturnCode string `json:"returnCode"`
	ReturnMsg  string `json:"returnMsg"`
	ResultCode string `json:"resultCode"`
	ErrCode    string `json:"errCode"`
	ErrCodeDes string `json:"errCodeDes"`
	Openid     string `json:"openid"`
	BankType   string `json:"bankType"`
	TotalFee   uint32 `json:"totalFee"`
	OutTradeNo string `json:"outTradeNo"`
}
