package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pushbackServer/database"
	"pushbackServer/push"
)

// SilentPushController 处理静默推送请求
// 通过deviceKey获取设备token并执行静默推送
func SilentPushController(c *gin.Context) {
	deviceKey := c.Param("deviceKey")
	if deviceKey == "" {
		c.JSON(http.StatusOK, failed(http.StatusBadRequest, "device key is empty"))
		return
	}

	if !database.DB.KeyExists(deviceKey) {
		c.JSON(http.StatusOK, failed(http.StatusBadRequest, "device key is not exist"))
		return
	}

	token, err := database.DB.DeviceTokenByKey(deviceKey)
	if err != nil {
		c.JSON(http.StatusOK, failed(http.StatusBadRequest, "failed to get device token: %v", err))
		return
	}

	if err := push.SilentPush(token); err != nil {
		c.JSON(http.StatusOK, failed(http.StatusInternalServerError, "failed to send silent push: %v", err))
		return
	}

	c.JSON(http.StatusOK, success())
}
