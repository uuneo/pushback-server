package database

import (
	"os"
	"path/filepath"
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

// getDataDir checks if /data directory exists and has write permission
// returns /data if available, otherwise returns ./
func getDataDir() string {
	dataDir := "/data"
	if _, err := os.Stat(dataDir); err == nil {
		// Check if we have write permission
		if err := os.MkdirAll(filepath.Join(dataDir, "test"), 0755); err == nil {
			_ = os.RemoveAll(filepath.Join(dataDir, "test"))
			return dataDir
		}
	}
	return "./"
}

func init() {
	if dsn := config.LocalConfig.System.Dsn; len(dsn) > 10 {
		if database, err := NewMySQL(dsn); err == nil {
			DB = database
			return
		}
	}

	DB = NewBboltdb(getDataDir())
}
