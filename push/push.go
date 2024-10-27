package push

import (
	"fmt"
	"github.com/uuneo/apns2"
	"github.com/uuneo/apns2/payload"
	"pushbackServer/config"
	"time"
)

// Push message to APNs server
func Push(params *config.ParamsResult, pushType apns2.EPushType) error {
	// 创建 payload，并填充通知标题、内容、声音和类别等字段
	pl := payload.NewPayload().
		AlertTitle(fmt.Sprint(params.Get(config.Title))).
		AlertSubtitle(fmt.Sprint(params.Get(config.Subtitle))).
		AlertBody(fmt.Sprint(params.Get(config.Body))).
		Sound(fmt.Sprint(params.Get(config.Sound))).
		TargetContentID(fmt.Sprint(params.Get(config.ID))).
		Category(fmt.Sprint(params.Get(config.Category)))

	// 添加自定义参数
	skipKeys := map[string]struct{}{
		config.DeviceKey:   {},
		config.DeviceToken: {},
		config.Title:       {},
		config.Body:        {},
		config.Sound:       {},
		config.Category:    {},
	}

	for pair := params.Params.Oldest(); pair != nil; pair = pair.Next() {
		if _, skip := skipKeys[pair.Key]; skip {
			continue
		}
		pl.Custom(pair.Key, pair.Value)
	}

	CLI := <-CLIENTS // 从池中获取一个客户端
	CLIENTS <- CLI   // 将客户端放回池中

	// 创建并发送通知
	resp, err := CLI.Push(&apns2.Notification{
		DeviceToken: fmt.Sprint(params.Get(config.DeviceToken)),
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

// SilentPush 静默推送
func SilentPush(deviceToken string) error {
	pl := payload.NewPayload().
		ContentAvailable().
		Custom("type", "token")

	CLI := <-CLIENTS // 从池中获取一个客户端
	CLIENTS <- CLI   // 将客户端放回池中

	resp, err := CLI.Push(&apns2.Notification{
		DeviceToken: deviceToken,
		Topic:       config.LocalConfig.Apple.Topic,
		Payload:     pl,
		Expiration:  time.Now().Add(24 * time.Hour),
		PushType:    apns2.PushTypeBackground,
	})

	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("APNs silent push failed: %s", resp.Reason)
	}
	return nil
}
