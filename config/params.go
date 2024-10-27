package config

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wk8/go-ordered-map/v2"
	"strings"
)

type ParamsResult struct {
	Params *orderedmap.OrderedMap[string, interface{}]
}

func NewParamsResult(c *gin.Context) *ParamsResult {
	main := &ParamsResult{
		Params: orderedmap.New[string, interface{}](),
	}
	main.HandlerParamsToMapOrder(c)
	main.SetDefault()
	return main
}

func (p *ParamsResult) NormalizeKey(s string) string {
	result := strings.TrimSpace(s)
	result = strings.ReplaceAll(result, "-", "")
	return strings.ToLower(result)
}

func (p *ParamsResult) Get(key string) interface{} {
	if value, ok := p.Params.Get(key); ok {
		return value
	}
	return ""
}

func (p *ParamsResult) SetDefault() {

	setDefault := func(params *orderedmap.OrderedMap[string, interface{}], key string, other func(string, interface{})) {
		newKey := p.NormalizeKey(key)
		if value, ok := params.Get(newKey); !ok || value == nil || len(fmt.Sprint(value)) == 0 {
			other(newKey, value)
		}
	}

	// 处理默认值
	setDefault(p.Params, AutoCopy, func(key string, value interface{}) {
		p.Params.Set(key, AutoCopyDefault)
	})

	setDefault(p.Params, Level, func(key string, value interface{}) {
		p.Params.Set(key, LevelDefault)
	})

	setDefault(p.Params, Body, func(key string, value interface{}) {
		p.Params.Set(key, "-No Content-")
	})

	setDefault(p.Params, Category, func(key string, value interface{}) {
		p.Params.Set(key, CategoryDefault)
	})

	setDefault(p.Params, ID, func(key string, value interface{}) {
		messageID, _ := uuid.NewUUID()
		p.Params.Set(key, messageID)
	})

	setDefault(p.Params, Sound, func(key string, value interface{}) {
		if value != nil && !strings.HasSuffix(fmt.Sprint(value), ".caf") {
			p.Params.Set(key, fmt.Sprint(value)+".caf")
		}
	})

}

func (p *ParamsResult) HandlerParamsToMapOrder(c *gin.Context) {

	result := orderedmap.New[string, interface{}]()
	result.Set(Host, c.Request.Host)
	// 兼容一下2.1.1之前的版本
	result.Set(Callback, c.Request.Host)

	switch len(c.Params) {
	case 1:
		result.Set(DeviceKey, c.Params[0].Value)
	case 2:
		result.Set(DeviceKey, c.Params[0].Value)
		result.Set(Body, c.Params[1].Value)
	case 3:
		result.Set(DeviceKey, c.Params[0].Value)
		result.Set(Title, c.Params[1].Value)
		result.Set(Body, c.Params[2].Value)
	case 4:
		result.Set(DeviceKey, c.Params[0].Value)
		result.Set(Title, c.Params[1].Value)
		result.Set(Subtitle, c.Params[2].Value)
		result.Set(Body, c.Params[3].Value)
	}

	// 获取所有url参数
	var params = c.Request.URL.Query()

	for k, v := range params {
		key := p.NormalizeKey(k)
		if value, ok := result.Get(key); !ok || value == "" {
			result.Set(key, v[0])
		}
	}

	if c.Request.Method == "POST" {
		contentType := c.Request.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "application/json") {
			var jsonData map[string]interface{}
			err := json.NewDecoder(c.Request.Body).Decode(&jsonData)
			if err == nil {
				for k, v := range jsonData {
					result.Set(p.NormalizeKey(k), v)
				}
			}
		} else {
			err := c.Request.ParseForm()
			if err == nil {
				for k, v := range c.Request.PostForm {
					result.Set(p.NormalizeKey(k), v)
				}
			}
		}

	}

	for pair := result.Oldest(); pair != nil; pair = pair.Next() {
		p.Params.Set(p.NormalizeKey(pair.Key), pair.Value)
	}

}
