package app

import (
	"context"
	"database/sql"
	"example/web-service-gin/internal/cache"
	"example/web-service-gin/internal/config"
	"example/web-service-gin/internal/handler"
	"example/web-service-gin/internal/metrics"
	"example/web-service-gin/internal/repository"
	"example/web-service-gin/internal/router"
	"example/web-service-gin/internal/storage"
	"example/web-service-gin/internal/tracing"
	"example/web-service-gin/internal/usecase"

	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
)

func Run(ctx context.Context) error {
	tp, err := tracing.InitTracer(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to init tracer")
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			slog.Error("Error shutting down tracer", slog.Any("error", err))
		}
	}()

	cfg, err := config.Init()
	if err != nil {
		return errors.Wrap(err, "Error initialization config")
	}

	//открываем sql.DB для миграций
	sqlDB, err := sql.Open("postgres", cfg.GetConnStr())
	if err != nil {
		return errors.Wrap(err, "failed to open sql.DB connection")
	}
	defer sqlDB.Close()

	// Запускаем миграции
	migrationsDir := "./database/migrations"
	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		return errors.Wrap(err, "failed to run migrations")
	}

	//Подклчаемся через pgx для основного кода
	conn, err := storage.GetConnect(cfg.GetConnStr())
	if err != nil {
		slog.Error("Unable to connect to database: %v\n", slog.Any("error", err))
		return errors.Wrap(err, "Error to connect to database")
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

	metrics.InitMetrics(cfg.Metrics.Port, cacheProvider)
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
