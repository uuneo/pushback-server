package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v3"
	"golang.org/x/net/context"
	"log"
	"os"
)

const formatConfigNameString = "The environment name you are using in gin mode is %s, and the path to the config is %s\n"

func init() {

	app := &cli.Command{
		Name: "PUSHBACK_SERVER",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Value:   ConfigFilePath,
				Usage:   "config file path",
				Aliases: []string{"c"},
				Sources: cli.EnvVars("PB_SERVER_CONFIG"),
			},
			&cli.StringFlag{
				Name:    "mode",
				Value:   "",
				Usage:   "server mode",
				Aliases: []string{"m"},
				Sources: cli.EnvVars("PB_SERVER_MODE"),
			},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			configPath := command.String("config")
			LocalVP = viper.New()
			LocalVP.SetConfigFile(configPath)
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

			if ServerMode := command.String("mode"); ServerMode != "" {
				if ServerMode == gin.DebugMode {
					LocalConfig.System.Mode = ServerMode
				} else {
					LocalConfig.System.Mode = gin.ReleaseMode
				}
			}

			log.Println("configTem file loaded successfully")
			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Println(err)
		return
	}

}
