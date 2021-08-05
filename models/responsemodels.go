package models

type HttpResponse struct {
	Code    int
	Message string
	Data    interface{}
}

func NewHttpResponse() HttpResponse {
	return HttpResponse{
		Code:    0,
		Message: "",
		Data:    make(map[string]string, 0),
	}
}
