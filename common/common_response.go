package common

type Response struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func NewCommonResponse(statusCode int, statusMsg string) Response {
	return Response{
		StatusCode: statusCode,
		StatusMsg:  statusMsg,
	}
}

type LoginResponse struct {
	Token string `json:"token"`
	Response
}
