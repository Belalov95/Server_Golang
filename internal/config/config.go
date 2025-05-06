package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		Port string
	}
}

func NewConfig() string {
	//загрузка переменных окружения из .env файла
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Ошибка при загрузке данных из .env файла: %v", slog.Any("error", err))
	}

	//Получение переменных окружения
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	//Формирование строки подключения к базе данных
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	return connStr
}
