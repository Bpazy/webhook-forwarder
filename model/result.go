package model

type Result[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func NewSuccessResult[T any](data T) Result[T] {
	return Result[T]{
		Success: true,
		Data:    data,
		Message: "",
	}
}

func NewFailedResult(message string) Result[any] {
	return Result[any]{
		Success: false,
		Data:    nil,
		Message: message,
	}
}
