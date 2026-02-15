package response

type ApiResponse[T any] struct {
	StatusCode int     `json:"statusCode"`
	Message    string  `json:"message"`
	Data       *T      `json:"data,omitempty"`
	Error      *string `json:"error,omitempty"`
}

func Success[T any](message string, data *T) ApiResponse[T] {
	return ApiResponse[T]{StatusCode: 200, Message: message, Data: data}
}

func Created[T any](message string, data *T) ApiResponse[T] {
	return ApiResponse[T]{StatusCode: 201, Message: message, Data: data}
}

func Removed[T any](message string, data *T) ApiResponse[T] {
	return ApiResponse[T]{StatusCode: 200, Message: message, Data: data}
}

func Error(statusCode int, message string, errs ...error) ApiResponse[any] {
	r := ApiResponse[any]{StatusCode: statusCode, Message: message}
	if len(errs) > 0 && errs[0] != nil {
		e := errs[0].Error()
		r.Error = &e
	}

	return r
}
