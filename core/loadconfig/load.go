package loadconfig

import (
	"github.com/spf13/viper"
	"sync"
)

var (
	Config     *config
	onceConfig sync.Once
)

type config struct {
	Mode         string `yaml:"mode"`
	Addr         string `yaml:"addr"`
	Name         string `yaml:"name"`
	URL          string `yaml:"url"`
	MaxPingCount int    `yaml:"max_ping_count"`
}

func BasicConfig() *config {
	onceConfig.Do(func() {
		Config = &config{
			Mode:         viper.GetString("mode"),
			Addr:         viper.GetString("addr"),
			Name:         viper.GetString("name"),
			URL:          viper.GetString("url"),
			MaxPingCount: viper.GetInt("max_ping_count"),
		}
	})
	return Config
}
