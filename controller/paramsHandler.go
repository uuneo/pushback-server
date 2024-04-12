package controller

import (
	"NewBearService/config"
	"NewBearService/database"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

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

	// 获取所有post参数
	if c.Request.Method == "POST" {
		var postParams = c.Request.PostForm
		for k, v := range postParams {
			key := config.UnifiedParameter(k)
			if value, ok := paramsResult[key]; !ok || value == "" {
				paramsResult[key] = v[0]
			} else {
				paramsResult[key] = v[0]
			}

		}
	}

	// 处理默认值
	if config.VerifyMap(paramsResult, config.IsArchive) == "" {
		paramsResult[config.IsArchive] = config.IsArchiveDefault
	}
	if config.VerifyMap(paramsResult, config.AutoCopy) == "" {
		paramsResult[config.AutoCopy] = config.AutoCopyDefault
	}
	if config.VerifyMap(paramsResult, config.Level) == "" {
		paramsResult[config.Level] = config.LevelA
	}
	if config.VerifyMap(paramsResult, config.Category) == "" {
		paramsResult[config.Category] = config.CategoryDefault
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
