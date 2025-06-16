package push

import (
	"fmt"
	"github.com/uuneo/apns2"
	"github.com/uuneo/apns2/payload"
	"pushbackServer/config"
	"sync"
	"time"
)

// Push message to APNs server
func Push(params *config.ParamsMap, pushType apns2.EPushType, token string) error {
	// 创建 payload，并填充通知标题、内容、声音和类别等字段
	pl := payload.NewPayload().
		AlertTitle(config.PMGet(params, config.Title)).
		AlertSubtitle(config.PMGet(params, config.Subtitle)).
		AlertBody(config.PMGet(params, config.Body)).
		Sound(config.PMGet(params, config.Sound)).
		TargetContentID(config.PMGet(params, config.ID)).
		ThreadID(config.PMGet(params, config.Group)).
		Category(config.PMGet(params, config.Category))

	// 添加自定义参数
	skipKeys := map[string]struct{}{
		config.DeviceKey:   {},
		config.DeviceToken: {},
		config.Title:       {},
		config.Body:        {},
		config.Sound:       {},
		config.Category:    {},
	}

	for pair := params.Oldest(); pair != nil; pair = pair.Next() {
		if _, skip := skipKeys[pair.Key]; skip {
			continue
		}
		pl.Custom(pair.Key, pair.Value)
	}

	CLI := <-CLIENTS // 从池中获取一个客户端
	CLIENTS <- CLI   // 将客户端放回池中

	// 创建并发送通知
	resp, err := CLI.Push(&apns2.Notification{
		DeviceToken: token,
		CollapseID:  fmt.Sprint(params.Value(config.ID)),
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

func MorePush(params *config.ParamsResult, pushType apns2.EPushType) error {
	if len(params.Results) > 0 {
		var (
			errors []error
			mu     sync.Mutex
			wg     sync.WaitGroup
		)
		for _, param := range params.Results {
			wg.Add(1)
			go func(p *config.ParamsMap) {
				defer wg.Done()
				if err := Push(p, pushType, params.DeviceToken); err != nil {
					fmt.Println(err.Error())
					mu.Lock()
					errors = append(errors, err)
					mu.Unlock()
				}
			}(param)
		}

		wg.Wait()
		if len(errors) > 0 {
			return fmt.Errorf("APNs push failed: %v", errors)
		}
	} else {

		if err := Push(params.Params, apns2.PushTypeAlert, params.DeviceToken); err != nil {
			return err
		}
	}
	return nil
}
