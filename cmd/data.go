package cmd

import (
	"strings"

	"github.com/spf13/viper"
)

const defaultPort = 6443

type Config struct {
	Environment []environment `yaml:"environment"`
}

type environment struct {
	Name        string `yaml:"name"`
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	Certificate string `yaml:"certificate"`
}

func ClusterName(name string) string {
	s := strings.Split(name, "/")
	last := s[len(s)-1]
	return last
}

func ReadToml() Config {
	var c Config
	viper.SetConfigName("gen")    // no need to include file extension
	viper.AddConfigPath("config") // set the path of your config file

	err := viper.ReadInConfig()
	if err != nil {
		errorf("Cannot Read toml file: %v", err)
	} else {
		c = Config{
			Environment: []environment{
				environment{
					Name:        "dev",
					Port:        viper.GetString("dev.port"),
					Host:        viper.GetString("dev.host"),
					Certificate: viper.GetString("dev.cert"),
				},
				environment{
					Name:        "stage",
					Port:        viper.GetString("stage.port"),
					Host:        viper.GetString("stage.host"),
					Certificate: viper.GetString("stage.cert"),
				},
				environment{
					Name:        "prod",
					Port:        viper.GetString("prod.port"),
					Host:        viper.GetString("prod.host"),
					Certificate: viper.GetString("prod.cert"),
				},
			},
		}
	}
	return c
}

func (d *Config) getHost(e string) string {
	for _, i := range d.Environment {
		if i.Name == e {
			return i.Host
		}
	}
	return ""
}

func (d *Config) getPort(e string) string {
	for _, i := range d.Environment {
		if i.Name == e {
			return i.Port
		}
	}
	return ""
}

func (d *Config) getCert(e string) string {
	for _, i := range d.Environment {
		if i.Name == e {
			return i.Certificate
		}
	}
	return ""
}
