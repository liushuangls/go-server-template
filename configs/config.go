package configs

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/liushuangls/go-server-template/pkg/xslog"
)

type Config struct {
	App   App          `yaml:"App"`
	DB    DB           `yaml:"DB"`
	Redis Redis        `yaml:"Redis"`
	Log   xslog.Config `yaml:"Log"`
	Jwt   Jwt          `yaml:"Jwt"`
}

type App struct {
	Addr string `yaml:"Addr"`
	Mode string `yaml:"Mode"`
}

type DB struct {
	Dialect     string `yaml:"Dialect"`
	DSN         string `yaml:"DSN"`
	MaxIdle     int    `yaml:"MaxIdle"`
	MaxActive   int    `yaml:"MaxActive"`
	MaxLifetime int    `yaml:"MaxLifetime"`
	AutoMigrate bool   `yaml:"AutoMigrate"`
}

type Redis struct {
	Addr     string `yaml:"Addr"`
	DB       int    `yaml:"DB"`
	Password string `yaml:"Password"`
}

type Jwt struct {
	Secret string `yaml:"Secret"`
	Issuer string `yaml:"Issuer"`
}

func (c *Config) IsDebugMode() bool {
	return c.App.Mode == "debug" || c.App.Mode == "local"
}

func (c *Config) IsReleaseMode() bool {
	return c.App.Mode == "release" || c.App.Mode == "prod"
}

func InitConfig() (*Config, error) {
	var cfg Config
	mode := os.Getenv("APP_MODE")
	if mode == "" {
		mode = "prod"
	}
	configPath := fmt.Sprintf("configs/%s.yaml", mode)
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if err := viper.ReadConfig(file); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalExact(&cfg); err != nil {
		return nil, err
	}
	return &cfg, err
}
