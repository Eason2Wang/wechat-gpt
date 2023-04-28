package entity

type PayCallbackRequest struct {
	ReturnCode string `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
	ResultCode string `json:"result_code"`
	ErrCode    string `json:"err_code"`
	ErrCodeDes string `json:"err_code_des"`
	Openid     string `json:"openid"`
	BankType   string `json:"bank_type"`
	TotalFee   uint32 `json:"total_fee"`
	OutTradeNo string `json:"out_trade_no"`
}
