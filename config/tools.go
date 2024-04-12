package config

import "strings"

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
