package controller

import (
	"fmt"
	"github.com/uuneo/apns2"
	"pushbackServer/config"
	"pushbackServer/push"
	"sync"
	"time"
)

func init() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			fmt.Println("开始检查未推送数据")
			NotPushedDataList.Range(func(key, value any) bool {
				data, ok := value.(*NotPushedData)
				if !ok {
					NotPushedDataList.Delete(key) // 类型异常也清除
					return true
				}

				now := time.Now()

				// 超过 24 小时未成功推送，直接清除
				if now.Sub(data.LastPushDate) > 24*time.Hour {
					NotPushedDataList.Delete(key)
					return true
				}

				// 推送节流策略：每次失败后等待 Count × 10 分钟
				nextTry := data.LastPushDate.Add(time.Duration(data.Count) * 10 * time.Minute)
				if nextTry.After(now) {
					return true // 还没到下一次推送时间，跳过
				}

				// 执行推送
				if err := push.Push(data.Params, data.PushType); err != nil {
					NotPushedDataList.Delete(key) // 推送失败直接删
				}

				return true
			})
		}
	}()

}

var NotPushedDataList sync.Map

type NotPushedData struct {
	ID           string               `json:"id"`
	CreateDate   time.Time            `json:"createDate"`
	LastPushDate time.Time            `json:"lastPushDate"`
	Count        int                  `json:"count"`
	Params       *config.ParamsResult `json:"params"`
	PushType     apns2.EPushType      `json:"pushType"`
}

// UpdateNotPushedData 更新已有记录，若不存在则添加
func UpdateNotPushedData(id string, params *config.ParamsResult, pushType apns2.EPushType) {
	if val, ok := NotPushedDataList.Load(id); ok {
		res := val.(*NotPushedData)
		res.LastPushDate = time.Now()
		res.Count++
		res.Params = params
		res.PushType = pushType
		NotPushedDataList.Store(id, data) // 可省略，但保持一致性
	} else {
		NotPushedDataList.Store(id, &NotPushedData{
			ID:           id,
			CreateDate:   time.Now(),
			LastPushDate: time.Now(),
			Count:        1,
			Params:       params,
			PushType:     pushType,
		})
	}
}

// RemoveNotPushedData 删除数据
func RemoveNotPushedData(id string) {
	NotPushedDataList.Delete(id)
}
