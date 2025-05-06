// основной файл программы, который запускает веб-сервер
package main

import (
	"context"
	"example/web-service-gin/internal/cache"
	"example/web-service-gin/internal/config"
	"example/web-service-gin/internal/handler"
	"example/web-service-gin/internal/repository"
	"example/web-service-gin/internal/router"
	"example/web-service-gin/internal/storage"
	"example/web-service-gin/internal/usecase"

	"fmt"
	"log"
	"log/slog"
	"os"
)

func main() {
	connStr := config.NewConfig()
	//присваиваем переменной Conn значение соединения
	//вызывается ф-ция GetConnect из пакета storage, которая устанавливает соединение с бд
	conn, err := storage.GetConnect(connStr)
	if err != nil {
		slog.Error("Unable to connect to database: %v\n", slog.Any("error", err))
	}
	//Эта строка использует ключевое слово defer, чтобы отложить выполнение функции Close до тех пор,
	// пока функция main не завершит выполнение. Это гарантирует, что соединение с базой данных будет закрыто,
	// даже если произойдет ошибка.
	defer conn.Close(context.Background())

	repo := repository.NewGoodsRepo(conn)

	cacheProvider := cache.New(repo)
	//Создает новый экземпляр бизнес логики (usecase) и передает ему соединение с бд
	uc := usecase.New(cacheProvider)

	//Создает новый обработчик HTTP-запросов (handler), передавая ему экземпляр бизнес логики (uc)
	//Обработчик будет использовать бизнес логику для обработки запросов от клиентов
	handle := handler.New(uc)

	//Создает новый роутер (маршрутизатор) для обработки HTTP-запросов
	router := router.GetRouter(handle)

	//Получаем переменную окружения
	appPort := os.Getenv("APP_PORT")

	//запуск сервера на указанном хосте и порту
	log.Printf("Запуск сервера на 0.0.0.0:%s\n", appPort)
	if err := router.Run(fmt.Sprintf(":%s", appPort)); err != nil {
		slog.Error("Ошибка при запуске сервера: %v", slog.Any("error", err))
	}
}
