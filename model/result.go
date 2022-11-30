package model

type Result struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewSuccessResult(data any) Result {
	return Result{
		Success: true,
		Data:    data,
		Message: "",
	}
}

func NewFailedResult(message string) Result {
	return Result{
		Success: false,
		Data:    nil,
		Message: message,
	}
}
