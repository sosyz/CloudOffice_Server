package log

import (
	"fmt"
	"time"
)

type ErrorInfo struct {
	code int
	msg  string
}

func NewError(code int, msg string) *ErrorInfo {
	// TODO: 存入数据库
	return &ErrorInfo{code, msg}
}

func (e *ErrorInfo) GetErrorCode() int {
	return e.code
}

func (e *ErrorInfo) GetErrorMsg() string {
	return e.msg
}

func Debug(tag, msg string) {
	fmt.Printf("%v [%s] %s\n", time.Now(), tag, msg)
}
