package config

import (
	"bytes"
	"fmt"
	"log/slog"

	_ "embed"

	"github.com/spf13/viper"
)

//go:embed config.yml
var defaultYMLFile []byte

type Config struct {
	App struct {
		Port string
	}
	Metrics struct {
		Port string
	}
	Db struct {
		Port     string
		Host     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	Cache struct {
		TTLSeconds int
	}
}

func (c *Config) GetConnStr() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		c.Db.User, c.Db.Password, c.Db.Host, c.Db.Port, c.Db.Name, c.Db.SSLMode)
}

func Init() (*Config, error) {
	// загружаем yml из embedded []byte
	viper.SetConfigType("yml")
	if err := viper.ReadConfig(bytes.NewBuffer(defaultYMLFile)); err != nil {
		slog.Error("Error reading embedded config", slog.Any("error", err))
		return nil, err
	}

	// ручное связывание из yml и env
	_ = viper.BindEnv("app.port", "APP_PORT")
	_ = viper.BindEnv("db.port", "DB_PORT")
	_ = viper.BindEnv("db.host", "DB_HOST")
	_ = viper.BindEnv("db.user", "DB_USER")
	_ = viper.BindEnv("db.password", "DB_PASSWORD")
	_ = viper.BindEnv("db.name", "DB_NAME")
	_ = viper.BindEnv("db.sslmode", "DB_SSLMODE")
	_ = viper.BindEnv("metrics.port", "METRICS_PORT")

	// Разбор конфигурации в структуру
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		slog.Error("Error filling data into structure", slog.Any("error", err))
		return nil, err
	}
	return &cfg, nil
}
