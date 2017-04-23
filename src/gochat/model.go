package gochat

type baseRequest struct {
	Sid      	string
	Skey       	string
	Uin      	string
	DeviceID	string
}

type Response struct {
	BaseResponse *baseResponse
}

type baseResponse struct {
	Ret    int
	ErrMsg string
}