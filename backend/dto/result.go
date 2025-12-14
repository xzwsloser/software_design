package dto

// 前后端交互结构体
type Result[T any] struct {
	// 响应是否成功
	Success 	bool 	`json:"success"`
	ErrorMsg	string 	`json:"error"`
	Data 		T 		`json:"data"`
}

func Ok() Result[string] {
	return Result[string] {
		Success: true,
		ErrorMsg: "",
		Data: "",
	}
}

func OkWithData[T any](data T) Result[T] {
	return Result[T] {
		Success: true,
		ErrorMsg: "",
		Data: data,
	}
}

func Fail(errMsg string) Result[string]{
	return Result[string] {
		Success: false,
		ErrorMsg: errMsg,
		Data: "",
	}
}

