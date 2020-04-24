package response

import (
	"github.com/lightserver-cn/golang-utils/sql"
)

type JsonResponse struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Json(code uint32, msg string, data interface{}, success bool) *JsonResponse {
	return &JsonResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func JsonData(data interface{}) *JsonResponse {
	return &JsonResponse{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
}

func JsonPageResult(result interface{}, page *sql.Pages) *JsonResponse {
	return JsonData(&sql.PageResult{
		Page:   page,
		Result: result,
	})
}

func JsonSuccess() *JsonResponse {
	return &JsonResponse{
		Code: 0,
		Msg:  "success",
		Data: nil,
	}
}

func JsonErrorCode(code uint32) *JsonResponse {
	return &JsonResponse{
		Code: code,
		Msg:  "failure",
		Data: nil,
	}
}

func JsonErrorMsg(message string) *JsonResponse {
	return &JsonResponse{
		Code: 400,
		Msg:  message,
		Data: nil,
	}
}

func JsonErrorCodeMsg(code uint32, msg string) *JsonResponse {
	return &JsonResponse{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

func JsonErrorData(code uint32, msg string, data interface{}) *JsonResponse {
	return &JsonResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

type Builder struct {
	Data map[string]interface{}
}

func NewEmptyResponseBuilder() *Builder {
	return &Builder{Data: make(map[string]interface{})}
}

func (b *Builder) Put(key string, value interface{}) *Builder {
	b.Data[key] = value
	return b
}

func (b *Builder) Build() map[string]interface{} {
	return b.Data
}

func (b *Builder) JsonResponse() *JsonResponse {
	return JsonData(b.Data)
}
