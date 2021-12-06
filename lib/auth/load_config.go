package auth

import (
	"sync"

	"github.com/spf13/viper"
)

var Config *config
var once sync.Once

type config struct {
	Admin         map[string]interface{} `yaml:"admin"`
	AdminUsername string
	AdminPassword string
	Users         map[string]interface{} `yaml:"users"`
	Path          string                 `yaml:"config_path"`
}

func LoadConfig() *config {
	once.Do(func() {
		Config = &config{
			Admin: viper.GetStringMap("auth.admin"),
			Users: viper.GetStringMap("auth.users"),
			Path:  viper.GetString("auth.config_path"),
		}
		for u, p := range Config.Admin {
			Config.AdminUsername = u
			Config.AdminPassword = p.(string)
			Config.Users[u] = p.(string)
		}
	})
	return Config
}
