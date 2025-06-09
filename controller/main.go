package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"github.com/uuneo/apns2"
	"net/http"
	"pushbackServer/config"
	"pushbackServer/database"
	"pushbackServer/push"
	"time"
)

// RegisterController 处理设备注册请求
// 支持 GET 和 POST 两种请求方式:
// GET: 检查设备key是否存在
// POST: 注册新的设备token
func RegisterController(c *gin.Context) {
	if c.Request.Method == "GET" {
		deviceKey := c.Param("deviceKey")
		if deviceKey == "" {
			c.JSON(http.StatusOK, failed(http.StatusBadRequest, "device key is empty"))
			return
		}
		if database.DB.KeyExists(deviceKey) {
			c.JSON(http.StatusOK, success())
			return
		} else {
			admin, ok := c.Get("admin")

			if ok && admin.(bool) {
				_, err := database.DB.SaveDeviceTokenByKey(deviceKey, "")
				if err != nil {
					c.JSON(http.StatusOK, failed(http.StatusBadRequest, "device key is not exist"))
					return
				}
				c.JSON(http.StatusOK, success())
				return
			}

			c.JSON(http.StatusOK, failed(http.StatusBadRequest, "device key is not exist"))
		}
		return
	}

	var err error
	var device DeviceInfo

	if err = c.BindJSON(&device); err != nil {
		c.JSON(http.StatusOK, failed(http.StatusBadRequest, "failed to get device token: %v", err))
		return
	}

	if device.Token == "" {
		c.JSON(http.StatusOK, failed(http.StatusBadRequest, "deviceToken is empty"))
		return
	}

	device.Key, err = database.DB.SaveDeviceTokenByKey(device.Key, device.Token)

	if err != nil {
		c.JSON(http.StatusOK, failed(http.StatusInternalServerError, "device registration failed: %v", err))
	}

	c.JSON(http.StatusOK, data(device))
}

// BaseController 处理基础推送请求
// 验证推送参数并执行推送操作
func BaseController(c *gin.Context) {
	ParamsResult := config.NewParamsResult(c)

	// 如果 title, subtitle 和 body 都为空，则返回错误
	if ParamsResult.IsNan {
		c.JSON(http.StatusOK, failed(http.StatusBadRequest, "title, subTitle, cipherText and body cannot be all empty"))
		return
	}

	if value, ok := ParamsResult.Params.Get(config.DeviceToken); !ok || len(fmt.Sprint(value)) <= 0 {
		deviceValue, deviceOk := ParamsResult.Params.Get(config.DeviceKey)
		if !deviceOk || len(fmt.Sprint(deviceValue)) <= 0 {
			c.JSON(http.StatusOK,
				failed(http.StatusBadRequest, "failed to get device token: deviceKey or deviceToken is required"))
			return
		}
		token, err := database.DB.DeviceTokenByKey(fmt.Sprint(deviceValue))
		if err != nil {
			c.JSON(http.StatusOK, failed(http.StatusBadRequest, "failed to get device token: %v", err))
			return
		}
		ParamsResult.Params.Set(config.DeviceToken, token)
	}

	err := push.Push(ParamsResult, apns2.PushTypeAlert)

	if err != nil {
		c.JSON(http.StatusOK, failed(http.StatusInternalServerError, "push failed: %v", err))
		return
	}

	// TODO: 这里需要判断是否是管理员
	admin, ok := c.Get("admin")
	// 如果是管理员，加入到未推送列表
	if ok && admin.(bool) {
		UpdateNotPushedData(ParamsResult.Get("id").(string), ParamsResult, apns2.PushTypeAlert)
	}

	c.JSON(http.StatusOK, success())
}

// GetPushToken 获取设备的推送token
// 通过deviceKey查询对应的推送token
func GetPushToken(c *gin.Context) {
	deviceKey := c.Param("deviceKey")
	fmt.Println(deviceKey)
	token, err := database.DB.DeviceTokenByKey(deviceKey)
	if err != nil {
		c.JSON(http.StatusOK, failed(http.StatusInternalServerError, "failed to get device token: %v", err))
		return
	}
	c.JSON(http.StatusOK, data(token))
}

// HomeController 处理首页请求
// 支持两种功能:
// 1. 通过id参数移除未推送数据
// 2. 生成二维码图片
func HomeController(c *gin.Context) {
	id := c.Query("id")
	if id != "" {
		RemoveNotPushedData(id)
		c.Status(http.StatusOK)
		return
	}

	url := "https://" + c.Request.Host

	code := c.Query("code")

	if code != "" {
		url = code
	}
	png, err := qrcode.Encode(url, qrcode.High, 1024)

	if err != nil {
		c.JSON(http.StatusOK, failed(http.StatusInternalServerError, "failed to generate QR code: %v", err))
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}

// Ping 处理心跳检测请求
// 返回服务器当前状态
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, CommonResp{
		Code:      http.StatusOK,
		Message:   "pong",
		Timestamp: time.Now().Unix(),
	})
}
