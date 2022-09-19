package settings

import "github.com/spf13/viper"

type Config struct {
	*viper.Viper
}

var settings *Config

func init() {
	Reset()
}

func Reset() {
	settings = &Config{viper.New()}
}

func SetConfigName(n string) {
	settings.SetConfigName(n)
}

func AddConfigPath(p string) {
	settings.AddConfigPath(p)
}

func ReadInConfig() error {
	return settings.ReadInConfig()
}

func GetString(k string) string {
	return settings.GetString(k)
}

func GetStringWithDefault(k string, d string) string {
	if !settings.IsSet(k) {
		return d
	}
	return settings.GetString(k)
}

func GetStringList(k string) []string {
	return settings.GetStringSlice(k)
}

func IsInStringList(k string, v string) bool {
	l := GetStringList(k)
	for _, s := range l {
		if s == v {
			return true
		}
	}

	return false
}

func GetBool(k string) bool {
	return settings.GetBool(k)
}

func GetBoolWithDefault(k string, d bool) bool {
	if !settings.IsSet(k) {
		return d
	}

	return settings.GetBool(k)
}
