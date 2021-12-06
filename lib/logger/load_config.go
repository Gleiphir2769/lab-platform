package logger

import (
	"sync"

	"github.com/spf13/viper"
)

var LogConfig *logConfig
var onceLog sync.Once

// LogConfig
type logConfig struct {
	Level      string `yaml:"level"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"maxsize"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
	BufferSize int    `yaml:"buffer_size"`
}

func LoadConfig() *logConfig {
	onceLog.Do(func() {
		LogConfig = &logConfig{
			Level:      viper.GetString("log.level"),
			Filename:   viper.GetString("log.filename"),
			MaxSize:    viper.GetInt("log.maxsize"),
			MaxAge:     viper.GetInt("log.max_age"),
			MaxBackups: viper.GetInt("log.max_backups"),
			BufferSize: viper.GetInt("log.buffer_size"),
		}
		GetLog()
	})
	return LogConfig
}
