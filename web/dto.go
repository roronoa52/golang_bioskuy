package web

type FormatResponsePaging struct {
	ResponseCode int         `json:"response-code"`
	Data         interface{} `json:"data"`
	Paging       Paging      `json:"paging"`
}

type Paging struct {
	Page      int `json:"page"`
	TotalData int `json:"total-data"`
}

type FormatResponse struct {
	ResponseCode int         `json:"response-code"`
	Data         interface{} `json:"data"`
}
