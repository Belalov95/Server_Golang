// основной файл программы, который запускает веб-сервер
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
	//вызывается ф-ция GetConnect из пакета storage, которая устанавливает соединение с бд
	conn, err := storage.GetConnect("postgresql://postgres:@postgres:5432/postgres")
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

	//запускаем сервер на localhost с портом 8080
	router.Run("0.0.0.0:8080")

}
