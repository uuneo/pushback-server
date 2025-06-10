package config

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wk8/go-ordered-map/v2"
	"strings"
)

// ParamsResult 结构体用于存储和管理请求参数
// 使用有序映射存储参数，保证参数的处理顺序
type ParamsResult struct {
	Params *orderedmap.OrderedMap[string, interface{}]
	IsNan  bool
}

// NewParamsResult 创建新的参数结果对象
// 参数:
//   - c: gin上下文对象，用于获取请求参数
//
// 返回:
//   - *ParamsResult: 初始化后的参数结果对象
func NewParamsResult(c *gin.Context) *ParamsResult {
	main := &ParamsResult{
		Params: orderedmap.New[string, interface{}](),
	}
	main.HandlerParamsToMapOrder(c)
	main.SetDefault()
	main.IsNan = ParamsNan(main)
	return main
}

// NormalizeKey 规范化参数键名
// 主要功能:
// 1. 去除首尾空格
// 2. 移除连字符
// 3. 转换为小写
// 参数:
//   - s: 需要规范化的键名字符串
//
// 返回:
//   - string: 规范化后的键名
func (p *ParamsResult) NormalizeKey(s string) string {
	result := strings.TrimSpace(s)
	result = strings.ReplaceAll(result, "-", "")
	result = strings.ReplaceAll(result, "_", "")
	return strings.ToLower(result)
}

// Get 获取参数值
// 参数:
//   - key: 参数键名
//
// 返回:
//   - interface{}: 参数值，如果不存在则返回空字符串
func (p *ParamsResult) Get(key string) interface{} {
	if value, ok := p.Params.Get(key); ok {
		return value
	}
	return ""
}

// SetDefault 设置参数的默认值
// 主要功能：
// 1. 为未设置或为空的参数设置默认值
// 2. 处理自动复制、消息级别、消息分类等参数的默认值
// 3. 为消息生成唯一ID
func (p *ParamsResult) SetDefault() {

	// setDefault 内部函数用于设置默认值
	// params: 参数映射
	// key: 需要设置默认值的键
	// other: 设置默认值的回调函数
	setDefault := func(params *orderedmap.OrderedMap[string, interface{}], key string, other func(string, interface{})) {
		newKey := p.NormalizeKey(key)
		if value, ok := params.Get(newKey); !ok || value == nil || len(fmt.Sprint(value)) == 0 {
			other(newKey, value)
		}
	}

	// 处理默认值
	// 设置自动复制功能的默认值
	setDefault(p.Params, AutoCopy, func(key string, value interface{}) {
		p.Params.Set(key, AutoCopyDefault)
	})

	// 设置消息级别的默认值
	setDefault(p.Params, Level, func(key string, value interface{}) {
		p.Params.Set(key, LevelDefault)
	})

	// 设置消息分类的默认值
	setDefault(p.Params, Category, func(key string, value interface{}) {
		p.Params.Set(key, CategoryDefault)
	})

	// 设置消息ID的默认值（使用UUID）
	setDefault(p.Params, ID, func(key string, value interface{}) {
		messageID, _ := uuid.NewUUID()
		p.Params.Set(key, messageID)
	})

}

// HandlerParamsToMapOrder 处理请求参数并转换为有序映射
// 主要功能：
// 1. 从URL路径参数中提取设备密钥、标题、副标题和内容
// 2. 从URL查询参数中获取额外参数
// 3. 处理POST请求的表单数据和JSON数据
// 4. 对参数进行便捷处理
// 5. 将处理后的参数保存到有序映射中
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
	convenientProcessor(result)

	for pair := result.Oldest(); pair != nil; pair = pair.Next() {
		p.Params.Set(p.NormalizeKey(pair.Key), pair.Value)
	}

}

// convenientProcessor 处理推送参数的便捷转换
// 主要功能：
// 1. 将 data/content/message/text 字段统一转换为 body
// 2. 处理 markdown 相关字段，设置对应的 category
// 3. 规范化 category 字段的值
// 4. 处理声音文件后缀
func convenientProcessor(params *orderedmap.OrderedMap[string, interface{}]) {
	// 如果没有 body 字段，尝试从其他字段转换
	if _, ok := params.Get(Body); !ok {
		if data, dataOk := params.Get(Data); dataOk {
			params.Set(Body, data)
			params.Delete(Data)
		} else if content, contentOk := params.Get(Content); contentOk {
			params.Set(Body, content)
			params.Delete(Content)
		} else if message, messageOk := params.Get(Message); messageOk {
			params.Set(Body, message)
			params.Delete(Message)
		} else if text, textOk := params.Get(Text); textOk {
			params.Set(Body, text)
			params.Delete(Text)
		}
	}

	// 处理 markdown 字段
	// 如果存在 markdown 字段，将其转换为 body 并设置 category 为 markdown
	if v, ok := params.Get(Markdown); ok {
		params.Set(Body, v)
		params.Set(Category, CategoryMarkdown)
		params.Delete(MD)
	}
	// 如果存在 md 字段，将其转换为 body 并设置 category 为 markdown
	if v, ok := params.Get(MD); ok {
		params.Set(Body, v)
		params.Set(Category, CategoryMarkdown)
		params.Delete(Markdown)
	}

	// 规范化 category 字段
	// 如果 category 不是默认值或 markdown，则设置为默认值
	if v, ok := params.Get(Category); ok {
		if v != CategoryDefault && v != CategoryMarkdown {
			params.Set(Category, CategoryDefault)
		}
	}

	// 处理声音文件后缀
	// 如果声音文件没有 .caf 后缀，则添加后缀
	if val, ok := params.Get(Sound); ok {
		if !strings.HasSuffix(val.(string), ".caf") {
			params.Set(Sound, val.(string)+".caf")
		}
	}
}

func ParamsNan(ParamsResult *ParamsResult) bool {
	var titleNan, subTitleNan, bodyNan, cipherNan = false, false, false, false
	title, titleOk := ParamsResult.Params.Get(Title)
	subTitle, subTitleOk := ParamsResult.Params.Get(Subtitle)
	body, bodyOk := ParamsResult.Params.Get(Body)
	cipherText, cipherTextOk := ParamsResult.Params.Get(CipherText)

	if !titleOk && !subTitleOk && !bodyOk && !cipherTextOk {
		return true
	}
	if titleOk {
		if title1 := strings.ReplaceAll(title.(string), " ", ""); len(title1) <= 0 {
			titleNan = true
		}
	}

	if subTitleOk {
		if subTitle1 := strings.ReplaceAll(subTitle.(string), " ", ""); len(subTitle1) <= 0 {
			subTitleNan = true
		}
	}

	if bodyOk {
		if body1 := strings.ReplaceAll(body.(string), " ", ""); len(body1) <= 0 {
			bodyNan = true
		}
	}

	if cipherTextOk {
		if cipherText1 := strings.ReplaceAll(cipherText.(string), " ", ""); len(cipherText1) <= 0 {
			cipherNan = true
		}
	}

	if bodyNan && !cipherNan {
		ParamsResult.Params.Set(Body, "--body is empty--")
	}

	return titleNan && subTitleNan && bodyNan && cipherNan
}
