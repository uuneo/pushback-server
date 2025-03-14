package push

import (
	"fmt"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"pushbackServer/config"
	"time"
)

// Push message to APNs server
func Push(params map[string]string, pushType apns2.EPushType) error {
	// 创建 payload，并填充通知标题、内容、声音和类别等字段
	pl := payload.NewPayload().
		AlertTitle(config.VerifyMap(params, config.Title)).
		AlertSubtitle(config.VerifyMap(params, config.Subtitle)).
		AlertBody(config.VerifyMap(params, config.Body)).
		Sound(config.VerifyMap(params, config.Sound)).
		Category(config.CategoryDefault)

	// 添加自定义参数
	skipKeys := map[string]struct{}{
		config.DeviceKey:   {},
		config.DeviceToken: {},
		config.Title:       {},
		config.Body:        {},
		config.Sound:       {},
	}

	for k, v := range params {
		k = config.UnifiedParameter(k)
		if _, skip := skipKeys[k]; skip {
			continue
		}
		fmt.Println("Custom parameter added:", k, v)
		pl.Custom(k, v)
	}

	// 设置通知组（线程 ID）
	if group := config.VerifyMap(params, config.Group); group != "" {
		pl = pl.ThreadID(group)
	}

	// 创建并发送通知
	resp, err := CLI.Push(&apns2.Notification{
		DeviceToken: params[config.DeviceToken],
		Topic:       config.LocalConfig.Apple.Topic,
		Payload:     pl.MutableContent(),
		Expiration:  time.Now().Add(24 * time.Hour),
		PushType:    pushType,
	})

	// 错误处理
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("APNs push failed: %s", resp.Reason)
	}
	return nil
}
