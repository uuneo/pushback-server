package push

import (
	"NewBearService/config"
	"fmt"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"time"
)

// Push  message to apns server
func Push(params map[string]string) error {

	pl := payload.NewPayload().
		AlertTitle(config.VerifyMap(params, config.Title)).
		AlertBody(config.VerifyMap(params, config.Body)).
		Sound(config.VerifyMap(params, config.Sound)).
		Category(config.VerifyMap(params, config.Category))

	for k, v := range params {
		k = config.UnifiedParameter(k)
		if k == config.DeviceKey ||
			k == config.DeviceToken ||
			k == config.Title ||
			k == config.Body ||
			k == config.Sound ||
			k == config.Category {
			continue
		}

		pl.Custom(k, v)

	}

	if group := config.VerifyMap(params, config.Group); group != "" {
		pl = pl.ThreadID(group)
	} else {
		pl = pl.ThreadID(config.DefaultGroup)
		params[config.Group] = config.DefaultGroup
	}

	resp, err := CLI.Push(&apns2.Notification{
		DeviceToken: params[config.DeviceToken],
		Topic:       config.LocalConfig.Apple.Topic,
		Payload:     pl.MutableContent(),
		Expiration:  time.Now().Add(24 * time.Hour),
	})

	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("APNS push failed: %s", resp.Reason)
	}
	return nil
}
