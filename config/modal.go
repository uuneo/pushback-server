package config

type Config struct {
	System System `mapstructure:"system" json:"system" yaml:"system"`
	Mysql  Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Apple  Apple  `mapstructure:"apple" json:"apple" yaml:"apple"`
}

type System struct {
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Post     string `mapstructure:"port" json:"port" yaml:"port"`
	Name     string `mapstructure:"name" json:"name" yaml:"name"`
	Mode     string `mapstructure:"mode" json:"mode" yaml:"mode"`
	DBType   string `mapstructure:"dbType" json:"dbType" yaml:"dbType"`
	DBPath   string `mapstructure:"dbPath" json:"dbPath" yaml:"dbPath"`
}

type Mysql struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     string `mapstructure:"port" json:"port" yaml:"port"`
	UserName string `mapstructure:"userName" json:"userName" yaml:"userName"`
	PassWord string `mapstructure:"passWord" json:"passWord" yaml:"passWord"`
}

type Apple struct {
	ApnsPrivateKey string `mapstructure:"apnsPrivateKey" json:"apnsPrivateKey" yaml:"apnsPrivateKey"`
	Topic          string `mapstructure:"topic" json:"topic" yaml:"topic"`
	KeyID          string `mapstructure:"keyID" json:"keyID" yaml:"keyID"`
	TeamID         string `mapstructure:"teamID" json:"teamID" yaml:"teamID"`
	Develop        bool   `mapstructure:"develop" json:"develop" yaml:"develop"`
}
