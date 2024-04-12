package config

import "github.com/spf13/viper"

// URL > QUERY > POST

const (
	LevelA = "active"
	LevelT = "timeSensitive"
	LevelP = "passive"
)

const (
	CategoryDefault  = "myNotificationCategory" // 模版标志
	AutoCopyDefault  = "0"                      // 默认自动复制
	IsArchiveDefault = "1"                      // 默认归档
	DeviceKey        = "devicekey"              // 设备key
	DeviceToken      = "devicetoken"            // 设备token
	Category         = "category"               // 类别
	Title            = "title"                  // 标题
	Body             = "body"                   // 内容
	IsArchive        = "isarchive"              // 是否归档
	Group            = "group"                  // 组
	DefaultGroup     = "Default"                // 默认组
	Sound            = "sound"                  // 声音
	AutoCopy         = "autocopy"               // 自动复制
	Level            = "level"                  // 等级
)

const (
	FilePathEnv     = "ALARM_PAW_CONFIG"
	DefaultFilePath = "./config/config.yaml"
	TestFilePath    = "/data/config.test.yaml"
	ReleaseFilePath = "/data/config.release.yaml"
)

var (
	LocalConfig *Config
	LocalVP     *viper.Viper
)
