package config

type Config struct {
	System System `mapstructure:"system" json:"system" yaml:"system"`
	Apple  Apple  `mapstructure:"apple" json:"apple" yaml:"apple"`
}

type System struct {
	User               string `mapstructure:"user" json:"user" yaml:"user"`
	Password           string `mapstructure:"password" json:"password" yaml:"password"`
	Address            string `mapstructure:"address" json:"address" yaml:"address"`
	Name               string `mapstructure:"name" json:"name" yaml:"name"`
	Debug              bool   `mapstructure:"debug" json:"debug" yaml:"debug"`
	Dsn                string `mapstructure:"dsn" json:"dsn" yaml:"dsn"`
	MaxApnsClientCount int    `mapstructure:"maxApnsClientCount" json:"maxApnsClientCount" yaml:"maxApnsClientCount"`
}

type Apple struct {
	ApnsPrivateKey string `mapstructure:"apnsPrivateKey" json:"apnsPrivateKey" yaml:"apnsPrivateKey"`
	Topic          string `mapstructure:"topic" json:"topic" yaml:"topic"`
	KeyID          string `mapstructure:"keyID" json:"keyID" yaml:"keyID"`
	TeamID         string `mapstructure:"teamID" json:"teamID" yaml:"teamID"`
	Develop        bool   `mapstructure:"develop" json:"develop" yaml:"develop"`
	AdminId        string `mapstructure:"adminId" json:"adminId" yaml:"adminId"`
}
