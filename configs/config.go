package configs

import (
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/liushuangls/go-server-template/pkg/xslog"
)

type Config struct {
	App   App          `koanf:"App"`
	DB    DB           `koanf:"DB"`
	Redis Redis        `koanf:"Redis"`
	Log   xslog.Config `koanf:"Log"`
	Jwt   Jwt          `koanf:"Jwt"`
}

type App struct {
	Addr string `koanf:"Addr"`
	Mode string `koanf:"Mode"`
}

type DB struct {
	Dialect     string `koanf:"Dialect"`
	DSN         string `koanf:"DSN"`
	MaxIdle     int    `koanf:"MaxIdle"`
	MaxActive   int    `koanf:"MaxActive"`
	MaxLifetime int    `koanf:"MaxLifetime"`
	AutoMigrate bool   `koanf:"AutoMigrate"`
}

type Redis struct {
	Addr     string `koanf:"Addr"`
	DB       int    `koanf:"DB"`
	Password string `koanf:"Password"`
}

type Jwt struct {
	Secret string `koanf:"Secret"`
	Issuer string `koanf:"Issuer"`
}

func (c *Config) IsDebugMode() bool {
	return c.App.Mode == "debug" || c.App.Mode == "local"
}

func (c *Config) IsReleaseMode() bool {
	return c.App.Mode == "release" || c.App.Mode == "prod"
}

func InitConfig() (*Config, error) {
	var (
		cfg Config
		err error
		k   = koanf.New(".")
	)
	mode := os.Getenv("APP_MODE")
	if mode == "" {
		mode = "prod"
	}
	configPath := fmt.Sprintf("configs/%s.yaml", mode)

	err = k.Load(file.Provider(configPath), yaml.Parser())
	if err != nil {
		return nil, fmt.Errorf("error loading file config: %v", err)
	}

	err = k.Load(env.Provider("_", env.Opt{
		Prefix: "",
		TransformFunc: func(k, v string) (string, any) {
			return k, v
		},
	}), nil)
	if err != nil {
		return nil, fmt.Errorf("error loading env config: %v", err)
	}

	err = k.Unmarshal("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal config: %v", err)
	}

	return &cfg, nil
}
