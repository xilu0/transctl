package cmd

const (
	baiduApi = "https://fanyi-api.baidu.com/api/trans/vip/translate"
	id       = "20221114001452792"
	secret   = "n8aROYumV2YSgQtFRkSu"
	salt     = 1435660288
)

type TransResult struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type responseBody struct {
	From        string        `json:"from"`         //": "en",
	To          string        `json:"to"`           // "to": "zh",
	TransResult []TransResult `json:"trans_result"` // "trans_result": [
	ErrorCode   string        `json:"error_code"`   //: "54001",
	ErrorMsg    string        `json:"error_msg"`    //: "Invalid Sign"
}
