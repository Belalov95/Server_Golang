package config

import (
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"app"`
	Db struct {
		Port     string `mapstructure:"port"`
		Host     string `mapstructure:"host"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"db"`
}

func Init() (*Config, error) {
	//загрузка переменных окружения из .env файла
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("Ошибка при загрузке данных из .env файла: %v", slog.Any("error", err))
	}

	viper.SetConfigName("config") //имя файла с yml
	viper.SetConfigType("yml")

	// Ищем в разных местах
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app/config")
	viper.AddConfigPath("/src/config")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./internal/config")

	//ручное связывание из yml и env
	_ = viper.BindEnv("app.port", "APP_PORT")
	_ = viper.BindEnv("db.port", "DB_PORT")
	_ = viper.BindEnv("db.host", "DB_HOST")
	_ = viper.BindEnv("db.user", "DB_USER")
	_ = viper.BindEnv("db.password", "DB_PASSWORD")
	_ = viper.BindEnv("db.name", "DB_NAME")

	//чтение конфигурационного файла config.yml
	if err := viper.ReadInConfig(); err != nil {
		slog.Error("Ошибка при чтении файла", slog.Any("error", err))
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		slog.Error("Ошибка при заполнении данных в структру", slog.Any("error", err))
	}
	return &cfg, nil
}
