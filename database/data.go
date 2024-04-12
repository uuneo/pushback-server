package database

var DB Database

// Database defines all the db operation
type Database interface {
	CountAll() (int, error)                                 //Get db records count
	DeviceTokenByKey(key string) (string, error)            //Get specified device's token
	SaveDeviceTokenByKey(key, token string) (string, error) //Create or update specified devices's token
	Close() error                                           //Close the database
}
