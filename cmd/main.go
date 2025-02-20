package main

import (
	"context"
	"example/web-service-gin/internal/handler"
	"example/web-service-gin/internal/router"
	"example/web-service-gin/internal/storage"
	"log"

	"github.com/shopspring/decimal"
)

// представление того, какие у меня будут данные
type footballstore struct {
	ID       string          `json:"id"`       //объявление номера // serial primary key
	Category string          `json:"category"` //категория товара: Одежда, обувь, аксессуары и тд.
	Name     string          `json:"name"`     //название товара:форма, бутсы, брелки и тд.
	Price    decimal.Decimal `json:"price"`    //цена товара
}

func main() {
	conn, err := storage.GetConnect("postgresql://postgres:@postgres:5432/postgres") //присваиваем переменной Conn значение соединения
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err) // %v\n позволяет включить конкретное сообщение об ошибке, без %v\n будет выведено просто сообщение которое в ""
	}
	defer conn.Close(context.Background()) //отложенное закрытие функции, т.е ф-ция которая идет после defer будет выполена после завершения основной ф-ции

	handle := handler.New(conn)
	router := router.GetRouter(handle)
	router.Run("0.0.0.0:8080") //запускаем сервер на localhost с портом 8080

}
