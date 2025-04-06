package config

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	c := loadDefaultConfig()
	if configBuffer, err := json.Marshal(c); err != nil {
		log.Println("Oops! Marshal config is failed. ", err)
		return nil, err
	} else if err = viper.ReadConfig(bytes.NewBuffer(configBuffer)); err != nil {
		log.Println("Oops! Read default config is failed. ", err)
		return nil, err
	}
	if err := viper.MergeInConfig(); err != nil {
		log.Println("Oops! Read config file is failed. ", err)
	}

	err := viper.Unmarshal(c)
	return c, err
}
