package database

import (
	"NewBearService/config"
	"database/sql"
	"fmt"
	"github.com/lithammer/shortuuid/v3"
	"log"
)

type MySQL struct {
}

var mysqlDB *sql.DB

func CreateDbSchema() string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS  `%s` (", config.LocalConfig.System.Name) +
		"    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT," +
		"    `key` VARCHAR(255) NOT NULL," +
		"    `token` VARCHAR(255) NOT NULL," +
		"    PRIMARY KEY (`id`)," +
		"    UNIQUE KEY `key` (`key`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
}

func NewMySQL(dsn string) Database {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatalf("failed to open database connection (%s): %v", dsn, err)
	}
	dbSchema := CreateDbSchema()
	_, err = db.Exec(dbSchema)
	if err != nil {
		log.Fatalf("failed to init database schema(%s): %v", dbSchema, err)
	}

	mysqlDB = db
	return &MySQL{}
}

func (d *MySQL) CountAll() (int, error) {
	var count int
	rawString := fmt.Sprintf("SELECT COUNT(1) FROM `%s`", config.LocalConfig.System.Name)
	err := mysqlDB.QueryRow(rawString).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (d *MySQL) DeviceTokenByKey(key string) (string, error) {
	var token string
	rawString := fmt.Sprintf("SELECT `token` FROM `%s` WHERE `key`=?", config.LocalConfig.System.Name)
	err := mysqlDB.QueryRow(rawString, key).Scan(&token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (d *MySQL) SaveDeviceTokenByKey(key, token string) (string, error) {
	if key == "" {
		// Generate a new UUID as the deviceKey when a new device register
		key = shortuuid.New()
	}
	rawString := fmt.Sprintf("INSERT INTO `%s` (`key`,`token`) VALUES (?,?) ON DUPLICATE KEY UPDATE `token`=?", config.LocalConfig.System.Name)

	_, err := mysqlDB.Exec(rawString, key, token, token)
	if err != nil {
		return "", err
	}

	return key, nil
}

func (d *MySQL) Close() error {
	return mysqlDB.Close()
}
