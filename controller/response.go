package controller

import (
	"fmt"
	"time"
)

type CommonResp struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// for the fast return success result
func success() CommonResp {
	return CommonResp{
		Code:      200,
		Message:   "success",
		Timestamp: time.Now().Unix(),
	}
}

// for the fast return failed result
func failed(code int, message string, args ...interface{}) CommonResp {
	return CommonResp{
		Code:      code,
		Message:   fmt.Sprintf(message, args...),
		Timestamp: time.Now().Unix(),
	}
}

// for the fast return result with custom data
func data(data interface{}) CommonResp {
	return CommonResp{
		Code:      200,
		Message:   "success",
		Timestamp: time.Now().Unix(),
		Data:      data,
	}
}
