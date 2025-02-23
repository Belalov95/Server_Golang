package main

import (
	"context"
	"example/web-service-gin/internal/handler"
	"example/web-service-gin/internal/router"
	"example/web-service-gin/internal/storage"
	"example/web-service-gin/internal/usecase"
	"log"
)

func main() {
	//присваиваем переменной Conn значение соединения
	conn, err := storage.GetConnect("postgresql://postgres:@postgres:5432/postgres")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	//отложенное закрытие функции, т.е ф-ция которая идет после defer будет выполена
	// после завершения основной ф-ции
	defer conn.Close(context.Background())

	uc := usecase.New(conn)
	handle := handler.New(uc)
	router := router.GetRouter(handle)

	//запускаем сервер на localhost с портом 8080
	router.Run("0.0.0.0:8080")

}
