package config

import (
	"strings"

	"github.com/spf13/viper"
)

// MongoDB configurations
type MongoDB struct {
	URL      string `json:"url" yaml:"url" toml:"url" mapstructure:"url"`
	Name     string `json:"name" yaml:"name" toml:"name" mapstructure:"name"`
	Login    string `json:"login" yaml:"login" toml:"login" mapstructure:"login"`
	Password string `json:"password" yaml:"password" toml:"password" mapstructure:"password"`
}

type RabbitMQ struct {
	URL string `json:"url" yaml:"url" toml:"url" mapstructure:"url"`
}

// Config of entire application.
type Config struct {
	MongoDB  *MongoDB  `json:"mongodb" yaml:"mongodb" toml:"mongodb" mapstructure:"mongodb"`
	RabbitMQ *RabbitMQ `json:"rabbitmq" yaml:"rabbitmq" toml:"rabbitmq" mapstructure:"rabbitmq"`
	Port     string    `json:"port" yaml:"port" toml:"port" mapstructure:"port"`
}

// NewMongoDB returns new default MongoDB configurations.
func NewMongoDB() (m *MongoDB) {
	m = new(MongoDB)
	m.URL = "mongodb://localhost:27017"
	m.Name = "pizza"
	return
}

// NewRabbitMQ returns new default RabbitMQ configurations.
func NewRabbitMQ() (r *RabbitMQ) {
	r = new(RabbitMQ)
	r.URL = "amqp://guest:guest@localhost:5672/"
	return
}

// New default configurations.
func New() (conf *Config) {
	conf = new(Config)
	conf.MongoDB = NewMongoDB()
	conf.RabbitMQ = NewRabbitMQ()
	return
}

// Load the Config from configuration files. This method panics on error.
func (c *Config) Load() *Config {
	var v = viper.New()
	v.SetConfigName("config") // } config.yaml
	v.SetConfigType("yaml")   // }
	v.AddConfigPath(".")      // for local use
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/manager") // for server
	v.SetEnvPrefix("MANAGER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	v.AutomaticEnv()

	if err := v.Unmarshal(c); err != nil {
		panic(err)
	}

	return c
}
