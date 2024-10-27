package database

import (
	"fmt"
	"pushbackServer/config"
)

var DB Database

// Database defines all the db operation
type Database interface {
	CountAll() (int, error)                                 //Get db records count
	DeviceTokenByKey(key string) (string, error)            //Get specified device's token
	SaveDeviceTokenByKey(key, token string) (string, error) //Create or update specified devices's token
	KeyExists(key string) bool
	Close() error //Close the database
}

func init() {
	switch config.LocalConfig.System.DBType {
	case "mysql":
		DB = NewMySQL(fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.LocalConfig.Mysql.Host,
			config.LocalConfig.Mysql.Port,
			config.LocalConfig.Mysql.Host,
			config.LocalConfig.Mysql.Port,
			config.LocalConfig.System.Name,
		))
	default:
		DB = NewBboltdb(config.LocalConfig.System.DBPath)

	}

}
