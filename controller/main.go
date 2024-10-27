package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sideshow/apns2"
	"net/http"
	"pushbackServer/config"
	"pushbackServer/database"
	"pushbackServer/push"
	"runtime"
	"strings"
	"time"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, CommonResp{
		Code:      200,
		Message:   "pong",
		Timestamp: time.Now().Unix(),
	})
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

func BaseController(c *gin.Context) {

	params, err := ToParamsHandler(c)
	if err != nil {
		c.JSON(http.StatusOK, failed(400, "failed to get device token: %v", err))
		return
	}
	err = push.Push(params, apns2.PushTypeAlert)

	if err != nil {
		c.JSON(http.StatusOK, failed(500, "push failed: %v", err))
		return
	}

	c.JSON(http.StatusOK, success())

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

	if c.Request.Method == "POST" {
		err := c.Request.ParseForm() // 确保解析表单数据
		if err == nil {
			for k, v := range c.Request.PostForm {
				paramsKey := config.UnifiedParameter(k)
				if paramsKey == config.DeviceKey && deviceKey == "" {
					deviceKey = v[0]
				} else if paramsKey == config.DeviceToken && deviceToken == "" {
					deviceToken = v[0]
				}
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
		"key":   newKey,
		"token": deviceToken,
	}))
}

func ToParamsHandler(c *gin.Context) (map[string]string, error) {
	var err error
	var paramsResult = make(map[string]string)
	// 获取所有路径参数
	switch len(c.Params) {

	case 1:
		paramsResult[config.DeviceKey] = c.Params[0].Value
	case 2:
		paramsResult[config.DeviceKey] = c.Params[0].Value
		paramsResult[config.Body] = c.Params[1].Value
	case 3:
		paramsResult[config.DeviceKey] = c.Params[0].Value
		paramsResult[config.Title] = c.Params[1].Value
		paramsResult[config.Body] = c.Params[2].Value
	case 4:
		paramsResult[config.DeviceKey] = c.Params[0].Value
		paramsResult[config.Category] = c.Params[1].Value
		paramsResult[config.Title] = c.Params[2].Value
		paramsResult[config.Body] = c.Params[3].Value
	}

	// 获取所有url参数
	var params = c.Request.URL.Query()

	for k, v := range params {
		key := config.UnifiedParameter(k)
		if value, ok := paramsResult[key]; !ok || value == "" {
			paramsResult[key] = v[0]
		}
	}

	if c.Request.Method == "POST" {
		err = c.Request.ParseForm()
		if err == nil {
			for k, v := range c.Request.PostForm {
				key := config.UnifiedParameter(k)
				paramsResult[key] = v[0] // 直接赋值
			}
		}

	}

	setDefault := func(params map[string]string, key, defaultValue string) {
		if config.VerifyMap(params, key) == "" {
			params[key] = defaultValue
		}
	}

	// 处理默认值
	setDefault(paramsResult, config.IsArchive, config.IsArchiveDefault)

	setDefault(paramsResult, config.AutoCopy, config.AutoCopyDefault)

	setDefault(paramsResult, config.Level, config.LevelA)

	setDefault(paramsResult, config.Category, config.CategoryDefault)

	setDefault(paramsResult, config.Body, "-No Content-")

	if config.VerifyMap(paramsResult, config.Sound) != "" && !strings.HasSuffix(paramsResult[config.Sound], ".caf") {
		paramsResult[config.Sound] = paramsResult[config.Sound] + ".caf"
	}

	if config.VerifyMap(paramsResult, config.DeviceToken) == "" {
		if config.VerifyMap(paramsResult, config.DeviceKey) == "" {
			err = errors.New("deviceKey or deviceToken is required")
			return nil, err
		}
		paramsResult[config.DeviceToken], err = database.DB.DeviceTokenByKey(paramsResult[config.DeviceKey])
		if err != nil {
			err = errors.New("failed to get device token: " + err.Error())
			return nil, err
		}
	}

	return paramsResult, nil
}

func ChangeKeyHandler(c *gin.Context) {

	oldKey := c.Query("oldKey")
	newKey := c.Query("newKey")
	deviceToken := c.Query("deviceToken")

	if deviceToken == "" {
		c.JSON(http.StatusOK, failed(400, "deviceToken is empty"))
		return
	}

	if database.DB.KeyExists(oldKey) && !database.DB.KeyExists(newKey) {
		_, err := database.DB.SaveDeviceTokenByKey(newKey, deviceToken)
		if err != nil {
			c.JSON(http.StatusOK, failed(500, "device registration failed: %v", err))
			return
		}
		c.JSON(http.StatusOK, data(map[string]string{
			"key":   newKey,
			"token": deviceToken,
		}))

	} else {
		c.JSON(http.StatusOK, failed(400, "deviceKey or newKey is invalid"))
	}

}
