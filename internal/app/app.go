package app

import (
	"context"
	"example/web-service-gin/database"
	"example/web-service-gin/internal/cache"
	"example/web-service-gin/internal/config"
	"example/web-service-gin/internal/handler"
	"example/web-service-gin/internal/metrics"
	"example/web-service-gin/internal/repository"
	"example/web-service-gin/internal/router"
	"example/web-service-gin/internal/storage"
	"example/web-service-gin/internal/tracing"
	"example/web-service-gin/internal/usecase"
	"time"

	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func Run(ctx context.Context) error {
	tp, err := tracing.Init(ctx)
	if err != nil {
		return errors.Wrap(err, "init tracer")
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			slog.Error("Error shutting down tracer", slog.Any("error", err))
		}
	}()

	cfg, err := config.Init()
	if err != nil {
		return errors.Wrap(err, "Init config")
	}

	if err := database.Migrate(cfg.GetConnStr()); err != nil {
		return errors.Wrap(err, "run migrations")
	}

	// Подклчаемся через pgx для основного кода
	pool, err := storage.GetConnect(ctx, cfg.GetConnStr())
	if err != nil {
		return errors.Wrap(err, "connect to database")
	}
	// Эта строка использует ключевое слово defer, чтобы отложить выполнение функции Close до тех пор,
	// пока функция main не завершит выполнение. Это гарантирует, что соединение с базой данных будет закрыто,
	// даже если произойдет ошибка.
	defer pool.Close()

	repo := repository.NewGoodsRepo(pool)

	//используем TTL из config
	cacheTTL := time.Duration(cfg.Cache.TTLSeconds) * time.Second
	cacheProvider := cache.New(repo, cacheTTL)
	// Создает новый экземпляр бизнес логики (usecase) и передает ему соединение с бд
	uc := usecase.New(cacheProvider)

	// Создает новый обработчик HTTP-запросов (handler), передавая ему экземпляр бизнес логики (uc)
	// Обработчик будет использовать бизнес логику для обработки запросов от клиентов
	handle := handler.New(uc)

	metrics.InitMetrics(cfg.Metrics.Port, cacheProvider)
	// Создает новый роутер (маршрутизатор) для обработки HTTP-запросов
	router := router.GetRouter(handle)

	// Получаем переменную окружения
	appPort := cfg.App.Port
	slog.Info("Starting the server", "address", "0.0.0.0:"+appPort)

	// запуск сервера на указанном хосте и порту
	if err := router.Run(fmt.Sprintf(":%s", appPort)); err != nil {
		return errors.Wrap(err, "server run")
	}
	return nil
}
