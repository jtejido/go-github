package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	*viper.Viper
	MaxLimit     int    `mapstructure:"max_limit"`
	UserLifetime int    `mapstructure:"user_lifetime"`
	Listen       string `mapstructure:"listen"`
	Debug        bool   `mapstructure:"debug"`
	Token        string `mapstructure:"token"`
}

func LoadConfig(filename string) (*Config, error) {
	conf := &Config{Viper: viper.New()}
	var configFile string

	if filename != "" {
		configFile = filename
	} else if os.Getenv("GITHUB_API_CONFIG") != "" {
		configFile = os.Getenv("GITHUB_API_CONFIG")
	} else {
		return nil, fmt.Errorf("Unable to find config file")
	}

	conf.Set("env", os.Getenv("env"))
	conf.SetConfigFile(configFile)
	if err := conf.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := conf.Unmarshal(conf); err != nil {
		return nil, err
	}

	go func() {
		conf.WatchConfig()
		// https://github.com/gohugoio/hugo/blob/master/watcher/batcher.go
		// https://github.com/spf13/viper/issues/609
		// for some reason this fires twice on a Win machine, and the way some editors save files.
		conf.OnConfigChange(func(e fsnotify.Event) {
			log.Println("Configuration has been changed...")
			// only re-read if file has been modified
			if err := conf.ReadInConfig(); err != nil {
				if err == nil {
					log.Println("Reading failed after configuration update: no data was read")
				} else {
					log.Fatalf("Reading failed after configuration update: %s \n", err.Error())
				}

				return
			} else {
				log.Println("Successfully re-read config file...")
			}

		})
	}()
	return conf, nil
}
