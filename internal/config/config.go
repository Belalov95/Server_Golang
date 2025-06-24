package config

import (
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"app"`
	Metrics struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"metrics"`
	Db struct {
		Port     string `mapstructure:"port"`
		Host     string `mapstructure:"host"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"db"`
}

func (c *Config) GetConnStr() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		c.Db.User, c.Db.Password, c.Db.Host, c.Db.Port, c.Db.Name, c.Db.SSLMode)
}

func Init() (*Config, error) {

	//Установка имени и типа конфигурационного файла
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	//Пути для поиска конфигурационного файла
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app/config")
	viper.AddConfigPath("/src/config")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./internal/config")

	//чтение конфигурационного файла config.yml
	if err := viper.ReadInConfig(); err != nil {
		slog.Error("Error read file", slog.Any("error", err))
	}

	//загрузка переменных окружения из .env файла
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("Error loading data from .env file:%v", slog.Any("error", err))
	}

	//ручное связывание из yml и env
	_ = viper.BindEnv("app.port", "APP_PORT")
	_ = viper.BindEnv("db.port", "DB_PORT")
	_ = viper.BindEnv("db.host", "DB_HOST")
	_ = viper.BindEnv("db.user", "DB_USER")
	_ = viper.BindEnv("db.password", "DB_PASSWORD")
	_ = viper.BindEnv("db.name", "DB_NAME")
	_ = viper.BindEnv("db.sslmode", "DB_SSLMODE")
	_ = viper.BindEnv("metrics.port", "METRICS_PORT")

	//Разбор конфигурации в структуру
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		slog.Error("Error filling data into structure", slog.Any("error", err))
	}
	return &cfg, nil
}
