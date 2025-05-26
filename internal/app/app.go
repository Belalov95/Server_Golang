package app

import (
	"context"
	"example/web-service-gin/internal/cache"
	"example/web-service-gin/internal/config"
	"example/web-service-gin/internal/handler"
	"example/web-service-gin/internal/metrics"
	"example/web-service-gin/internal/repository"
	"example/web-service-gin/internal/router"
	"example/web-service-gin/internal/storage"
	"example/web-service-gin/internal/usecase"

	"fmt"
	"log/slog"
	"os"
)

func Run() error {
	cfg, err := config.Init()
	if err != nil {
		slog.Error("Ошибка при инициализации конфига", slog.Any("error", err))
	}
	//присваиваем переменной Conn значение соединения
	//вызывается ф-ция GetConnect из пакета storage, которая устанавливает соединение с бд
	//Формирование строки подключения к базе данных
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.Name)
	conn, err := storage.GetConnect(connStr)
	if err != nil {
		slog.Error("Unable to connect to database: %v\n", slog.Any("error", err))
		os.Exit(1)
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

	metrics.InitMetrics("8082", cacheProvider)
	//Создает новый роутер (маршрутизатор) для обработки HTTP-запросов
	router := router.GetRouter(handle)

	//Получаем переменную окружения
	appPort := cfg.App.Port
	slog.Info("Starting the server", "address", "0.0.0.0:"+appPort)

	//запуск сервера на указанном хосте и порту
	if err := router.Run(fmt.Sprintf(":%s", appPort)); err != nil {
		slog.Error("Error when server starting: %v", slog.Any("error", err))
	}
	return nil
}
