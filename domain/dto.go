package domain

type GetOriginalReq struct {
	Short string
}

type GetInfoReq struct {
	Short string
}

type SaveReq struct {
	Original string
}

type HttpResponse struct {
	StatusCode int         `json:"status_code"`
	StatusDesc string      `json:"status_desc"`
	Error      string      `json:"error"`
	Payload    interface{} `json:"payload"`
}
