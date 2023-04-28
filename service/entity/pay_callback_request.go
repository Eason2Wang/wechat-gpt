package entity

type PayCallbackRequest struct {
	ReturnCode    string `json:"returnCode"`
	ResultCode    string `json:"resultCode"`
	BankType      string `json:"bankType"`
	TotalFee      uint32 `json:"totalFee"`
	OutTradeNo    string `json:"outTradeNo"`
	TransactionId string `json:"transactionId"`
}
