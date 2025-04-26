// основной файл программы, который запускает веб-сервер
package main

import (
	"context"
	"example/web-service-gin/internal/handler"
	"example/web-service-gin/internal/router"
	"example/web-service-gin/internal/storage"
	"example/web-service-gin/internal/usecase"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	//загрузка переменных окружения из .env файла
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Ошибка при загрузке данных из .env файла: %v", err)
	}

	//Получение переменных окружения
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	//Формирование строки подключения к базе данных
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	//присваиваем переменной Conn значение соединения
	//вызывается ф-ция GetConnect из пакета storage, которая устанавливает соединение с бд
	conn, err := storage.GetConnect(connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	//Эта строка использует ключевое слово defer, чтобы отложить выполнение функции Close до тех пор,
	// пока функция main не завершит выполнение. Это гарантирует, что соединение с базой данных будет закрыто,
	// даже если произойдет ошибка.
	defer conn.Close(context.Background())

	//Создает новый экземпляр бизнес логики (usecase) и передает ему соединение с бд
	uc := usecase.New(conn)

	//Создает новый обработчик HTTP-запросов (handler), передавая ему экземпляр бизнес логики (uc)
	//Обработчик будет использовать бизнес логику для обработки запросов от клиентов
	handle := handler.New(uc)

	//Создает новый роутер (маршрутизатор) для обработки HTTP-запросов
	router := router.GetRouter(handle)

	//Получаем переменную окружения
	appPort := os.Getenv("APP_PORT")

	//запуск сервера на указанном хосте и порту
	log.Printf("Запуск сервера на 0.0.0.0:%s\n", appPort)
	err = router.Run("0.0.0.0" + appPort)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
