package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/uuneo/apns2"
	"net/http"
	"pushbackServer/config"
	"pushbackServer/database"
	"pushbackServer/push"
	"runtime"
	"strings"
	"time"
)

func RegisterController(c *gin.Context) {

	if c.Request.Method == "GET" {
		deviceKey := c.Param("device_key")
		if deviceKey == "" {
			c.JSON(http.StatusOK, failed(400, "device key is empty"))
			return
		}
		if database.DB.KeyExists(deviceKey) {
			c.JSON(http.StatusOK, success())
			return
		} else {
			admin, ok := c.Get("admin")

			if ok && admin.(string) == config.LocalConfig.Apple.AdminId {
				_, err := database.DB.SaveDeviceTokenByKey(deviceKey, "")
				if err != nil {
					c.JSON(http.StatusOK, failed(400, "device key is not exist"))
					return
				}
				c.JSON(http.StatusOK, success())
				return
			}

			c.JSON(http.StatusOK, failed(400, "device key is not exist"))
		}
		return
	}

	var err error
	var device DeviceInfo

	if err = c.BindJSON(&device); err != nil {
		c.JSON(http.StatusOK, failed(400, "failed to get device token: %v", err))
		return
	}

	if device.Token == "" {
		c.JSON(http.StatusOK, failed(400, "deviceToken is empty"))
		return
	}

	device.Key, err = database.DB.SaveDeviceTokenByKey(device.Key, device.Token)

	if err != nil {
		c.JSON(http.StatusOK, failed(500, "device registration failed: %v", err))
	}

	c.JSON(http.StatusOK, data(device))
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

func GetPushToken(c *gin.Context) {
	deviceKey := c.Param("deviceKey")
	fmt.Println(deviceKey)
	token, err := database.DB.DeviceTokenByKey(deviceKey)
	if err != nil {
		c.JSON(http.StatusOK, failed(500, "failed to get device token: %v", err))
		return
	}
	c.JSON(http.StatusOK, data(token))
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
		paramsResult[config.Title] = c.Params[1].Value
		paramsResult[config.Subtitle] = c.Params[2].Value
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
		contentType := c.Request.Header.Get("Content-Type")
		if contentType == "application/json" {
			var jsonData map[string]interface{}
			err = json.NewDecoder(c.Request.Body).Decode(&jsonData)

			if err == nil {
				for k, v := range jsonData {
					key := config.UnifiedParameter(k)
					paramsResult[key] = fmt.Sprintf("%v", v) // 转换为字符串存储
				}
			}
		} else {
			err = c.Request.ParseForm()
			if err == nil {
				for k, v := range c.Request.PostForm {
					key := config.UnifiedParameter(k)
					paramsResult[key] = v[0] // 直接赋值
				}
			}
		}

	}

	setDefault := func(params map[string]string, key, defaultValue string) {
		if config.VerifyMap(params, key) == "" {
			params[key] = defaultValue
		}
	}

	// 处理默认值
	setDefault(paramsResult, config.AutoCopy, config.AutoCopyDefault)

	setDefault(paramsResult, config.Level, config.LevelDefault)

	setDefault(paramsResult, config.Body, "-No Content-")

	setDefault(paramsResult, config.Category, config.CategoryDefault)

	if config.VerifyMap(paramsResult, config.ID) == "" {
		messageID, _ := uuid.NewUUID()
		paramsResult[config.ID] = messageID.String()
	}

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

func QRCode(c *gin.Context) {
	url := config.LocalConfig.System.HostName
	var png []byte
	png, err := qrcode.Encode(url, qrcode.Medium, 1024)
	if err != nil {
		c.JSON(200, "生成二维码失败")
		return
	}
	_, _ = c.Writer.Write(png)
}

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
		"version": "0.2.0",
		"build":   "",
		"arch":    runtime.GOOS + "/" + runtime.GOARCH,
		"commit":  "",
		"devices": devices,
	})

}
