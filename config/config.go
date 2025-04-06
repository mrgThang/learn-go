package config

import (
	"fmt"
	"net/url"
)

type Config struct {
	Mysql           MysqlConfig `json:"mysql" mapstructure:"mysql"`
	MigrationFolder string      `json:"migration_folder" mapstructure:"migration_folder"`
}

type MysqlConfig struct {
	Host     string `json:"host" mapstructure:"host"`
	Port     string `json:"port" mapstructure:"port"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
	Database string `json:"database" mapstructure:"database"`
	Options  string `json:"options" mapstructure:"options" yaml:"options"`
}

func (c MysqlConfig) DSN() string {
	options := c.Options
	if options != "" {
		if options[0] != '?' {
			options = "?" + options
		}
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s",
		c.Username,
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.Database,
		options)
}

func (c MysqlConfig) String() string {
	return fmt.Sprintf("mysql://%s", c.DSN())
}

func loadDefaultConfig() *Config {
	return &Config{
		Mysql: MysqlConfig{
			Host:     "localhost",
			Port:     "3307",
			Username: "root",
			Password: "secret",
			Database: "learn_go",
			Options:  "?parseTime=true",
		},
		MigrationFolder: "file://migrations",
	}
}
