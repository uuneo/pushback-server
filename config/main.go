package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v3"
	"golang.org/x/net/context"
	"log"
	"os"
)

// init 初始化函数
// 用于初始化配置和命令行参数
// 主要功能:
// 1. 设置命令行参数
// 2. 读取配置文件
// 3. 监听配置文件变化
// 4. 设置服务器运行模式
func init() {

	app := &cli.Command{
		Name: "PUSHBACK_SERVER",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "addr",
				Usage:   "Server listen address",
				Sources: cli.EnvVars("PB_SERVER_ADDR"),
				Value:   "",
			},
			&cli.StringFlag{
				Name:    "config",
				Value:   "/data/config.yaml",
				Usage:   "config file path",
				Aliases: []string{"c"},
				Sources: cli.EnvVars("PB_SERVER_CONFIG"),
			},
			&cli.StringFlag{
				Name:    "dsn",
				Usage:   "MySQL DSN user:pass@tcp(host)/dbname",
				Value:   "",
				Sources: cli.EnvVars("PB_SERVER_DSN"),
			},
			&cli.Int8Flag{
				Name:    "maxApnsClientCount",
				Value:   0,
				Usage:   "max apns client count, 0 means no limit",
				Aliases: []string{"max"},
				Sources: cli.EnvVars("PB_MAX_APNS_CLIENT_COUNT"),
			},
			&cli.BoolFlag{
				Name:    "debug",
				Value:   false,
				Usage:   "enable debug mode",
				Sources: cli.EnvVars("PB_DEBUG"),
			},
			&cli.BoolFlag{
				Name:    "develop",
				Value:   false,
				Usage:   "enable push dev mode",
				Aliases: []string{"dev"},
				Sources: cli.EnvVars("PB_DEVELOP"),
			},
			&cli.StringFlag{
				Name:    "user",
				Value:   "",
				Usage:   "server user",
				Aliases: []string{"u"},
				Sources: cli.EnvVars("PB_USER"),
			},
			&cli.StringFlag{
				Name:    "password",
				Value:   "",
				Usage:   "server password",
				Aliases: []string{"p"},
				Sources: cli.EnvVars("PB_PASSWORD"),
			},
		},
		Authors: []any{"neo@uuneo.com"},
		Action: func(ctx context.Context, command *cli.Command) error {
			LocalVP = viper.New()

			configPath := command.String("config")

			// Check if config file exists at specified path
			if _, err := os.Stat(configPath); err != nil {
				// If not found, try ./config.yaml
				configPath = ""
				if _, err = os.Stat("./config.yaml"); err == nil {
					configPath = "./config.yaml"
				}
			}

			if configPath != "" {

				LocalVP.SetConfigFile(configPath)
				LocalVP.SetConfigType("yaml")
				err := LocalVP.ReadInConfig()
				if err != nil {
					LocalConfig = DefaultConfig
					log.Printf("load default config successfully, error: %v", err)
				} else {
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
					log.Println("configTem file loaded successfully")
				}

			} else {
				LocalConfig = DefaultConfig
				log.Println("config file not found, use default config")
			}

			if command.Bool("develop") {
				LocalConfig.Apple.Develop = true
			}

			if command.Bool("debug") {
				LocalConfig.System.Debug = true
			}

			if command.String("user") != "" {
				LocalConfig.System.User = command.String("user")
			}
			if password := command.String("password"); password != "" {
				LocalConfig.System.Password = password
			}
			if address := command.String("addr"); address != "" {
				LocalConfig.System.Address = address
			}
			if maxCount := command.Int8("maxApnsClientCount"); maxCount > 0 {
				LocalConfig.System.MaxApnsClientCount = int(maxCount)
			}
			if dsn := command.String("dsn"); len(dsn) > 10 {
				LocalConfig.System.Dsn = dsn
			}

			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Println(err)
		return
	}

}
