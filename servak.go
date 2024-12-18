package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// представление того, какие у меня будут данные
type footballstore struct {
	ID       string  `json:"id"`       //объявление номера // serial primary key
	Category string  `json:"category"` //категория товара: Одежда, обувь, аксессуары и тд.
	Name     string  `json:"name"`     //название товара:форма, бутсы, брелки и тд.
	Price    float64 `json:"price"`    //цена товара // 1.12 -> 112 -> bigint (в golang используй decimal либу)
}

// TODO: убрать это в БД, создать таблицу и положить скрипт sql в ./migration/init.sql и убрать var goods = []footballstore
// CREATE TABLE store ...
// объявляю изначальные тестовые данные, список товаров
// var goods = []footballstore{
// 	{ID: "1", Category: "Одежда", Name: "Футбольная форма", Price: 25.99},
// 	{ID: "2", Category: "Обувь", Name: "Бутсы", Price: 32.99},
// 	{ID: "3", Category: "Аксессуары", Name: "Брелок в виде мяча", Price: 8.99},
// }
// func Connect() sql.Conn
// gorm db.Exec().Where(``).Params(``), pgx или sql -> db.Exec(`select * from store`)

// реализовать func listStore, updateStore, getStoreByID, insertStore

// var globalConn sql.Conn

// назначаем функцию обработчика пути к конечной точке
func main() {
	// TODO:
	// Example: postgres://user:password@postgres_name:service_port/db?sslmode=disable
	// conn, err := Connect(postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable)
	// globalConn = conn

	router := gin.Default()                //создаем маршрутизатор, который помогает серверу направлять входящие запросы к нужной ф-ции
	router.GET("/goods", listGoods)        //определяем маршрут "/goods" и указываем ф-цию
	router.GET("/goods/:id", getGoodsByID) //определение маршрута, чтобы искать по id
	router.POST("/goods", postGoods)       //определение маршрута для POST-запроса на адрес /goods

	router.Run("localhost:8080") //запускаем сервер на localhost с портом 8080
}

// получаем список товаров
func listGoods(c *gin.Context) {
	// goods, err= := listStore()
	c.IndentedJSON(http.StatusOK, goods)
}

// postGoods добавляет товары используя данные в формате JSON, которые приходят в теле HTTP-запроса
func postGoods(c *gin.Context) {
	var newGoods footballstore

	//функция BindJSON используется для привязки(заполнения) данных из полученного JSON к переменной newGoods
	if err := c.BindJSON(&newGoods); err != nil { //проверка на ошибку
		// TODO: отдать ошибку
		// если есть ошибка после getStoreByID + "error": err.Error()
		// c.IndentedJSON(http.StatusNotFound, gin.H{"error": "err.Error()"}) //отправляем ответ лиенту что товар не найден
		return
	}
	//добавялем новые товары в магазин
	// goods = append(goods, newGoods)              //добавляем newGoods в массив goods
	// err := insertStore(newGoods)
	c.IndentedJSON(http.StatusCreated, newGoods) //оправляем клиенту ответ в формате JSON
}

// getGoodsByID ищет альбом, у которого ID совпадает с уникальным идентификатором id (например, 1,2...)
// Если такой товар найден, то функция возвращает его как ответ
func getGoodsByID(c *gin.Context) {
	id := c.Param("id") //получаем id из параметров URL

	// TODO: good, err := getStoreByID(id)
	// //просматриваем список товаров и ищем в нем товар с нужным id, который запросил клиент
	// for _, a := range goods { //цикл, который перебирает каждый товар в списке товаров
	// 	if a.ID == id { //проверка, совпадает ли ID текущего альбома с id, который был получен из запроса
	// 		c.IndentedJSON(http.StatusOK, a)
	// 		return
	// 	}
	// }

	// TODO: пиши на англ
	// если есть ошибка после getStoreByID + "error": err.Error()
	c.IndentedJSON(http.StatusNotFound, gin.H{"сообщение": "товар не был найден"}) //отправляем ответ лиенту что товар не найден
}

//func getGoodsByID(c *gin.Context) {
//	id:=c.Param("id")
//	var good footballstore
//  if err:=db.First(&good, id).Error; err!=nil{
//  c.IntentedJSON(http.StatusNotFound, gin.H{"error":"this product not found"})
//  return
//  }
//  c.IntentedJSON(http.StatusOK, good)
//}
