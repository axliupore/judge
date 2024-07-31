package response

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type JudgeResponse struct {
	Output string `json:"output"`
	Time   int64  `json:"time"`
	Memory int64  `json:"memory"`
	FileId string `json:"fileId"`
}
