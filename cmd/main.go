package main

import (
	"context"
	"example/web-service-gin/internal/handler"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

// представление того, какие у меня будут данные
type footballstore struct {
	ID       string          `json:"id"`       //объявление номера // serial primary key
	Category string          `json:"category"` //категория товара: Одежда, обувь, аксессуары и тд.
	Name     string          `json:"name"`     //название товара:форма, бутсы, брелки и тд.
	Price    decimal.Decimal `json:"price"`    //цена товара
}

func getConnect(connStr string) (*pgx.Conn, error) { //создается для соединения с бд
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	return conn, nil
}

func main() {
	var err error                                                             //эта переменная хранит в себе данные об ошибках
	conn, err := getConnect("postgresql://postgres:@localhost:5433/postgres") //присваиваем переменной Conn значение соединения
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err) // %v\n позволяет включить конкретное сообщение об ошибке, без %v\n будет выведено просто сообщение которое в ""
	}
	defer conn.Close(context.Background()) //отложенное закрытие функции, т.е ф-ция которая идет после defer будет выполена после завершения основной ф-ции

	handle := handler.New(conn)
	router := getRouter(handle)

	router.Run("0.0.0.0:8080") //запускаем сервер на localhost с портом 8080

}
func getRouter(handle *handler.Handle) *gin.Engine { // ф-ция, которая создает и настраивает HTTP-маршрутизатор

	router := gin.Default()                        //создаем маршрутизатор, который помогает серверу направлять входящие запросы к нужной ф-ции
	router.GET("/goods", handle.ListStore)         //определяем маршрут "/goods" и указываем ф-цию
	router.GET("/goods/:id", handle.GetGoodByID)   //определение маршрута, чтобы искать по id
	router.POST("/goods", handle.InsertStore)      //определение маршрута для POST-запроса на адрес /goods
	router.PUT("/goods/:id", handle.UpdateStore)   //определение маршрута для обновления данных о товаре
	router.DELETE("/goods/:id", handle.DeleteById) //определение маршрута для удаления

	return router
}
