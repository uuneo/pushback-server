package config

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const formatConfigNameString = "The environment name you are using in gin mode is %s, and the path to the config is %s\n"

func init() {
	var configTem string
	flag.StringVar(&configTem, "c", "", "choose configTem file.")
	flag.Parse()
	if configTem == "" { // 判断命令行参数是否为空
		if configEnv := os.Getenv(FilePathEnv); configEnv == "" { // 判断 internal.FilePathEnv 常量存储的环境变量是否为空
			switch gin.Mode() {
			case gin.DebugMode:
				configTem = DefaultFilePath
				fmt.Printf(formatConfigNameString, gin.DebugMode, configTem)
			case gin.ReleaseMode:
				configTem = ReleaseFilePath
				fmt.Printf(formatConfigNameString, gin.ReleaseMode, configTem)
			case gin.TestMode:
				configTem = TestFilePath
				fmt.Printf(formatConfigNameString, gin.TestMode, configTem)
			}
		} else { // internal.FilePathEnv 常量存储的环境变量不为空 将值赋值于config
			configTem = configEnv
			fmt.Printf(formatConfigNameString, FilePathEnv, configTem)
		}
	} else { // 命令行参数不为空 将值赋值于config
		fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%s\n", configTem)
	}

	LocalVP = viper.New()
	LocalVP.SetConfigFile(configTem)
	LocalVP.SetConfigType("yaml")
	err := LocalVP.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error configTem file: %s \n", err))
	}
	LocalVP.WatchConfig()

	LocalVP.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("configTem file changed:", e.Name)
		if err = LocalVP.Unmarshal(&LocalConfig); err != nil {
			fmt.Println(err)
		}
	})
	if err = LocalVP.Unmarshal(&LocalConfig); err != nil {
		fmt.Println(err)
	}

}

func GetDsn() string {
	return LocalConfig.Mysql.Host + ":" +
		LocalConfig.Mysql.Port +
		"@tcp(" +
		LocalConfig.Mysql.Host + ":" +
		LocalConfig.Mysql.Port + ")/" +
		LocalConfig.System.Name +
		"?charset=utf8mb4&parseTime=True&loc=Local"
}

func VerifyMap(data map[string]string, key string) string {
	if value, ok := data[key]; ok {
		return value
	}
	return ""
}

func UnifiedParameter(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}
