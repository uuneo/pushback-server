package database

import (
	"fmt"
	"github.com/lithammer/shortuuid/v3"
	"go.etcd.io/bbolt"
	"log"
	"os"
	"path/filepath"
	"pushbackServer/config"
	"sync"
)

// BboltDB implement Database interface with ETCD's bbolt
type BboltDB struct {
}

var dbOnce sync.Once
var BBDB *bbolt.DB

func NewBboltdb(dataDir string) Database {
	bboltSetup(dataDir)
	return &BboltDB{}
}

func (d *BboltDB) CountAll() (int, error) {
	var keypairCount int
	err := BBDB.View(func(tx *bbolt.Tx) error {
		keypairCount = tx.Bucket([]byte(config.LocalConfig.System.Name)).Stats().KeyN
		return nil
	})

	if err != nil {
		return 0, err
	}

	return keypairCount, nil
}

func (d *BboltDB) Close() error {
	return BBDB.Close()
}

func (d *BboltDB) DeviceTokenByKey(key string) (string, error) {
	var token string
	err := BBDB.View(func(tx *bbolt.Tx) error {
		if bs := tx.Bucket([]byte(config.LocalConfig.System.Name)).Get([]byte(key)); bs == nil {
			return fmt.Errorf("failed to get [%s] device token from database", key)
		} else {
			token = string(bs)
			return nil
		}
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

// SaveDeviceToken create or update device token of specified key

func (d *BboltDB) SaveDeviceTokenByKey(key, deviceToken string) (string, error) {
	err := BBDB.Update(func(tx *bbolt.Tx) error {

		bucket := tx.Bucket([]byte(config.LocalConfig.System.Name))
		// If the deviceKey is empty or the corresponding deviceToken cannot be obtained from the database,
		// it is considered as a new device registration
		if key == "" {
			// Generate a new UUID as the deviceKey when a new device register
			key = shortuuid.New()
		}
		// update the deviceToken
		return bucket.Put([]byte(key), []byte(deviceToken))
	})

	if err != nil {
		return "", err
	}

	return key, nil
}

// bboltSetup set up the bbolt database
func bboltSetup(dataDir string) {
	dbOnce.Do(func() {
		log.Printf("init database [%s]...", dataDir)
		if _, err := os.Stat(dataDir); os.IsNotExist(err) {
			if err = os.MkdirAll(dataDir, 0755); err != nil {
				log.Fatalf("failed to create database storage dir(%s): %v", dataDir, err)
			}
		} else if err != nil {
			log.Fatalf("failed to open database storage dir(%s): %v", dataDir, err)
		}

		bboltDB, err := bbolt.Open(filepath.Join(dataDir, config.LocalConfig.System.Name+".db"), 0600, nil)
		if err != nil {
			log.Fatalf("failed to create database file(%s): %v", filepath.Join(dataDir, config.LocalConfig.System.Name+".db"), err)
		}
		err = bboltDB.Update(func(tx *bbolt.Tx) error {
			_, err = tx.CreateBucketIfNotExists([]byte(config.LocalConfig.System.Name))
			return err
		})
		if err != nil {
			log.Fatalf("failed to create database bucket: %v", err)
		}

		BBDB = bboltDB
	})
}

// KeyExists 检查指定的 key 是否存在于数据库中，只返回 bool 值
func (d *BboltDB) KeyExists(key string) bool {
	err := BBDB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(config.LocalConfig.System.Name))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", config.LocalConfig.System.Name)
		}
		// 检查 key 是否存在
		if bucket.Get([]byte(key)) != nil {
			return nil // key 存在，返回 nil 表示没有错误
		}
		return fmt.Errorf("key not found") // key 不存在，返回错误
	})

	// 如果 err 为 nil，说明 key 存在，否则 key 不存在
	return err == nil
}
