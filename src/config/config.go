package config

import (
	"os"

	"github.com/spf13/viper"
)

type ConfigMode string

const (
	ConfigModeDevelopment ConfigMode = "development"
	ConfigModeRelease     ConfigMode = "release"
)

type Server struct {
	Host string
	Port int32
}

type Postgres struct {
	Host     string
	Port     int32
	Username string
	Password string
	Database string
}

type Config struct {
	Server   Server
	Postgres Postgres
}

type ConfigBuilder struct {
	config Config
	loader *viper.Viper
	err    error
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		loader: viper.New(),
	}
}

func (c *ConfigBuilder) Load() *ConfigBuilder {
	mode := ConfigMode(os.Getenv("GIN_MODE"))

	if mode == ConfigModeRelease {
		c.loader.AutomaticEnv()
	}

	baseDir, err := os.Getwd()
	if err != nil {
		c.err = err
		return c
	}

	c.loader.AddConfigPath(baseDir)
	c.loader.SetConfigName(".env")
	c.loader.SetConfigType("env")

	if err := c.loader.ReadInConfig(); err != nil {
		c.err = nil
		return c
	}

	c.config = Config{
		Server: Server{
			Host: c.loader.GetString("SERVER_HOST"),
			Port: c.loader.GetInt32("SERVER_PORT"),
		},
		Postgres: Postgres{
			Host:     c.loader.GetString("POSTGRES_HOST"),
			Port:     c.loader.GetInt32("POSTGRES_PORT"),
			Username: c.loader.GetString("POSTGRES_USERNAME"),
			Password: c.loader.GetString("POSTGRES_PASSWORD"),
			Database: c.loader.GetString("POSTGRES_DATABASE"),
		},
	}

	return c
}

func (c *ConfigBuilder) ServerHost(host string) *ConfigBuilder {
	c.config.Server.Host = host
	return c
}

func (c *ConfigBuilder) ServerPort(port int32) *ConfigBuilder {
	c.config.Server.Port = port
	return c
}

func (c *ConfigBuilder) PostgresHost(host string) *ConfigBuilder {
	c.config.Postgres.Host = host
	return c
}

func (c *ConfigBuilder) PostgresPort(port int32) *ConfigBuilder {
	c.config.Postgres.Port = port
	return c
}

func (c *ConfigBuilder) PostgresUsername(username string) *ConfigBuilder {
	c.config.Postgres.Username = username
	return c
}

func (c *ConfigBuilder) PostgresPassword(password string) *ConfigBuilder {
	c.config.Postgres.Password = password
	return c
}

func (c *ConfigBuilder) PostgresDatabase(database string) *ConfigBuilder {
	c.config.Postgres.Database = database
	return c
}

func (c *ConfigBuilder) Build() (*Config, error) {
	if c.err != nil {
		return nil, c.err
	}

	return &c.config, nil
}
