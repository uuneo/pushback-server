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

			if ok && admin.(string) == config.LocalConfig.Apple.AdminId {
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

func BaseController(c *gin.Context) {

	ParamsResult := config.NewParamsResult(c)

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
	if ok && admin.(string) == config.LocalConfig.Apple.AdminId {
		UpdateNotPushedData(ParamsResult.Get("id").(string), ParamsResult, apns2.PushTypeAlert)
	}

	c.JSON(http.StatusOK, success())

}

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

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, CommonResp{
		Code:      http.StatusOK,
		Message:   "pong",
		Timestamp: time.Now().Unix(),
	})
}
