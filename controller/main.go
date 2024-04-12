package controller

import (
	"NewBearService/config"
	"NewBearService/database"
	"NewBearService/push"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"time"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, CommonResp{
		Code:      200,
		Message:   "pong",
		Timestamp: time.Now().Unix(),
	})
}

func BaseController(c *gin.Context) {

	params, err := ToParamsHandler(c)
	if err != nil {
		c.JSON(http.StatusOK, failed(400, "failed to get device token: %v", err))
		return
	}
	err = push.Push(params)

	if err != nil {
		c.JSON(http.StatusOK, failed(500, "push failed: %v", err))
		return
	}

	c.JSON(http.StatusOK, success())

}

func GetInfo(c *gin.Context) {
	devices, _ := database.DB.CountAll()
	c.JSON(200, map[string]interface{}{
		"version": "1.0.0",
		"build":   "",
		"arch":    runtime.GOOS + "/" + runtime.GOARCH,
		"commit":  "",
		"devices": devices,
	})

}

func RegisterController(c *gin.Context) {
	var deviceKey, deviceToken string

	for _, v := range c.Params {
		paramsKey := config.UnifiedParameter(v.Key)
		if paramsKey == config.DeviceKey {
			deviceKey = v.Value
		} else if paramsKey == config.DeviceToken {
			deviceToken = v.Value
		}
	}

	for k, v := range c.Request.URL.Query() {
		paramsKey := config.UnifiedParameter(k)
		if paramsKey == config.DeviceKey && deviceKey == "" {
			deviceKey = v[0]
		} else if paramsKey == config.DeviceToken && deviceToken == "" {
			deviceToken = v[0]
		}
	}

	if c.Request.Method == "POST" {
		for k, v := range c.Request.PostForm {
			paramsKey := config.UnifiedParameter(k)
			if paramsKey == config.DeviceKey && deviceKey == "" {
				deviceKey = v[0]
			} else if paramsKey == config.DeviceToken && deviceToken == "" {
				deviceToken = v[0]
			}
		}
	}

	if deviceToken == "" {
		c.JSON(http.StatusOK, failed(400, "deviceToken is empty"))
		return
	}

	newKey, err := database.DB.SaveDeviceTokenByKey(deviceKey, deviceToken)
	if err != nil {
		c.JSON(http.StatusOK, failed(500, "device registration failed: %v", err))
	}

	c.JSON(http.StatusOK, data(map[string]string{
		"key":          newKey,
		"device_key":   newKey,
		"device_token": deviceToken,
	}))
}

func RegisterPush(c *gin.Context) {
	phone := c.Query("phone")
	code := c.Query("code")
	key := c.Query("key")

	if phone == "" {
		c.JSON(200, gin.H{
			"success": false,
			"msg":     "参数错误",
		})
		return
	} else if phone != "" && code == "" {
		c.JSON(200, gin.H{
			"success": true,
			"msg":     "发送成功",
		})
	} else if phone != "" || code != "" || key != "" {

		c.JSON(200, gin.H{
			"success": true,
			"msg":     "注册成功",
		})
		return
	} else {
		c.JSON(200, gin.H{
			"success": false,
			"msg":     "参数错误",
		})
		return
	}

}

func CallbackController(c *gin.Context) {
	// TODO 回调处理
	id := c.Params[0].Value
	fmt.Println("callback id: ", id)
	c.JSON(200, true)
}
