package http_handlers

import (
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/bufbuild/connect-go"
)

type HTTPResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func NewHTTPResponse(data interface{}) *HTTPResponse {
	resp := &HTTPResponse{}

	switch data.(type) {
	case e.ResponseError:
		resp.Code = int(data.(e.ResponseError).Code())
		resp.Msg = data.(e.ResponseError).Error()
	case error:
		resp.Code = int(connect.CodeInternal)
		resp.Msg = data.(error).Error()
	default:
		resp.Msg = "Success"
		resp.Data = data
	}

	return resp
}
