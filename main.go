package main

import (
	"context"
	"log"
	"net/http"

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

var Conn *pgx.Conn //эта переменная хранит в себе соединение с бд

func main() {
	var err error                                                                                   //эта переменная хранит в себе данные об ошибках
	Conn, err = pgx.Connect(context.Background(), "postgresql://postgres:@localhost:5433/postgres") //присваиваем переменной Conn значение соединения
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err) // %v\n позволяет включить конкретное сообщение об ошибке, без %v\n будет выведено просто сообщение которое в ""
	}
	defer Conn.Close(context.Background()) //отложенное закрытие функции, т.е ф-ция которая идет после defer будет выполена после завершения основной ф-ции

	router := gin.Default()                 //создаем маршрутизатор, который помогает серверу направлять входящие запросы к нужной ф-ции
	router.GET("/goods", listStore)         //определяем маршрут "/goods" и указываем ф-цию
	router.GET("/goods/:id", getGoodByID)   //определение маршрута, чтобы искать по id
	router.POST("/goods", insertStore)      //определение маршрута для POST-запроса на адрес /goods
	router.PUT("/goods/:id", updateStore)   //определение маршрута для обновления данных о товаре
	router.DELETE("/goods/:id", DeleteById) //определение маршрута для удаления

	router.Run("0.0.0.0:8080") //запускаем сервер на localhost с портом 8080
}

// создается для получения данных из таблицы бд
func listStore(c *gin.Context) {
	rows, err := Conn.Query(context.Background(), `SELECT "id", "category", "name", "price" from "footballstore"`) //выполнение SQL запроса через соединение с бд
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when executing sql query"}) //ошибка при выполнении sql запроса
		return
	}
	defer rows.Close()

	var goods []footballstore //создается срез

	for rows.Next() { //метод, который используется для перебора строк
		var good footballstore                                              //объявляем переменную good, чтобы хранить данные для каждой строки, которую мы считваем из бд
		err := rows.Scan(&good.ID, &good.Category, &good.Name, &good.Price) //извлекаем данные из текущей струтуры базы данны(id, catogory ...) и присваиваем их полям структуры(&id, &category ...), чтобы можно было работать с этими данными
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when reading data"}) //Ошибка при чтении данных
			return
		}
		goods = append(goods, good) //добавляем извлеченные товары в наш срез
	}
	if err := rows.Err(); err != nil { //проверка на ошибку после перебора
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error after sorting through the items"}) //Ошибка после перербора товаров
		return
	}
	c.JSON(http.StatusOK, goods) //отправляем список товаров в формате JSON
}

// создаем эту ф-цию для обновления данных в магазине (в базе данных)
func updateStore(c *gin.Context) {
	id := c.Param("id") //получяем id нужного товара для обновления

	var updatedGoods footballstore //создаем переменную updatedGoods чтобы хранить в ней обновленные товары

	if err := c.BindJSON(&updatedGoods); err != nil { //считываем с помощью BindJSON новые данные которые отправил клиент и передаем их переменной updatedGoods
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error unconnecting data"}) //некоректные данные
		return
	}
	query := "UPDATE footballstore SET category = $1, name = $2, price = $3 WHERE id = $4"                             //обновляем данные через SQL
	_, err := Conn.Exec(context.Background(), query, updatedGoods.Category, updatedGoods.Name, updatedGoods.Price, id) //выполняем SQL запрос для обновления данных и присваиваем новые значения переменным
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error during the update"}) //ошибка при обнолении
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "The good has been successfully updated"}) //отправляем клиенту ответ об успешном обновлении
}

// создается для поиска товара по его id
func getGoodByID(c *gin.Context) {
	id := c.Param("id") //создаем для поиска товара по id

	var good footballstore //создаем переменную для хранения данных о продукте

	query := "SELECT id, name, category, price FROM footballstore WHERE id = $1" // выполняем SQL запрос, где выдается конкретный id, в данном случае 1
	row := Conn.QueryRow(context.Background(), query, id)                        //используется чтобы выдать только 1 строку, в данном случае  id строку

	err := row.Scan(&good.ID, &good.Category, &good.Name, &good.Price) //Метод Scan извлекает значения из результата запроса и присваивает их полям структуры good.
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Good not found"}) //товар не найден
		return
	}
	c.JSON(http.StatusOK, good)
}

// добавляем новую запись в бд
func insertStore(c *gin.Context) {
	var newGood footballstore

	if err := c.BindJSON(&newGood); err != nil { //считываем json данные и присваиваем их переменной newGood
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "couldn't assign data"})
		return
	}
	query := "INSERT INTO footballstore (category, name, price) VALUES ($1, $2,  $3)"
	_, err := Conn.Exec(context.Background(), query, newGood.Category, newGood.Name, newGood.Price)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't assign data"}) //не удалось присвоить данные
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "The good added successfully"})
}
func DeleteById(c *gin.Context) {
	id := c.Param("id")
	query := "DELETE FROM footballstore WHERE id = $1"
	_, err := Conn.Exec(context.Background(), query, id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "id not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "The good has been deleated"})
}
