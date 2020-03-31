package response

type JsonResponse struct {
	ErrorCode int         `json:"error_code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Success   bool        `json:"success"`
}

func Json(errorCode int, message string, data interface{}, success bool) *JsonResponse {
	return &JsonResponse{
		ErrorCode: errorCode,
		Message:   message,
		Data:      data,
		Success:   success,
	}
}

func JsonData(data interface{}) *JsonResponse {
	return &JsonResponse{
		ErrorCode: 0,
		Data:      data,
		Success:   true,
		Message:   "success",
	}
}

func JsonSuccess() *JsonResponse {
	return &JsonResponse{
		ErrorCode: 0,
		Message:   "success",
		Data:      nil,
		Success:   true,
	}
}

func JsonErrorCode(code int) *JsonResponse {
	return &JsonResponse{
		ErrorCode: code,
		Message:   "failure",
		Data:      nil,
		Success:   false,
	}
}

func JsonErrorMsg(message string) *JsonResponse {
	return &JsonResponse{
		ErrorCode: 0,
		Message:   message,
		Data:      nil,
		Success:   false,
	}
}

func JsonErrorCodeMsg(code int, message string) *JsonResponse {
	return &JsonResponse{
		ErrorCode: code,
		Message:   message,
		Data:      nil,
		Success:   false,
	}
}

func JsonErrorData(code int, message string, data interface{}) *JsonResponse {
	return &JsonResponse{
		ErrorCode: code,
		Message:   message,
		Data:      data,
		Success:   false,
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
