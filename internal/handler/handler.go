package handler

import (
	"example/web-service-gin/internal/usecase"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

// представление того, какие у меня будут данные
type footballstore struct {
	ID       string          `json:"id"`       //объявление номера // serial primary key
	Category string          `json:"category"` //категория товара: Одежда, обувь, аксессуары и тд.
	Name     string          `json:"name"`     //название товара:форма, бутсы, брелки и тд.
	Price    decimal.Decimal `json:"price"`    //цена товара
}
type Handle struct {
	goodsUsecase *usecase.GoodsUsecase
}

func New(goodsUsecase *usecase.GoodsUsecase) *Handle {
	return &Handle{goodsUsecase: goodsUsecase}
}

// создается для получения данных из таблицы бд
func (h *Handle) ListStore(c *gin.Context) {
	goods, err := h.goodsUsecase.ListStore(c) //присваиваем переменной goods список товаров
	if err != nil {                           //проверка на ошибку после перебора
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error after sorting through the items"}) //Ошибка после перербора товаров
		return
	}
	c.JSON(http.StatusOK, goods) //отправляем список товаров в формате JSON
}

// создаем эту ф-цию для обновления данных в магазине (в базе данных)
func (h *Handle) UpdateStore(c *gin.Context) {
	id := c.Param("id") //получяем id нужного товара для обновления

	var updatedGoods footballstore //создаем переменную updatedGoods чтобы хранить в ней обновленные товары

	if err := c.BindJSON(&updatedGoods); err != nil { //считываем с помощью BindJSON новые данные которые отправил клиент и передаем их переменной updatedGoods
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error unconnecting data"}) //некоректные данные
		return
	}
	h.goodsUsecase.UpdateStore(c, id, updatedGoods)
	c.JSON(http.StatusOK, gin.H{"message": "The good has been successfully updated"}) //отправляем клиенту ответ об успешном обновлении
}

// создается для поиска товара по его id
func (h *Handle) GetGoodByID(c *gin.Context) {
	id := c.Param("id") //создаем для поиска товара по id

	var good footballstore //создаем переменную для хранения данных о продукте
	/*
		query := "SELECT id, name, category, price FROM footballstore WHERE id = $1" // выполняем SQL запрос, где выдается конкретный id, в данном случае 1
		row := h.conn.QueryRow(context.Background(), query, id)                      //используется чтобы выдать только 1 строку, в данном случае  id строку

		err := row.Scan(&good.ID, &good.Category, &good.Name, &good.Price) //Метод Scan извлекает значения из результата запроса и присваивает их полям структуры good.
	*/
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Good not found"}) //товар не найден
		return
	}
	c.JSON(http.StatusOK, good)
}

// добавляем новую запись в бд
func (h *Handle) InsertStore(c *gin.Context) {
	var newGood footballstore

	if err := c.BindJSON(&newGood); err != nil { //считываем json данные и присваиваем их переменной newGood
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "couldn't assign data"})
		return
	}
	/*
		query := "INSERT INTO footballstore (category, name, price) VALUES ($1, $2,  $3)"
		_, err := h.conn.Exec(context.Background(), query, newGood.Category, newGood.Name, newGood.Price)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't assign data"}) //не удалось присвоить данные
			return
		}
	*/
	c.JSON(http.StatusOK, gin.H{"message": "The good added successfully"})
}
func (h *Handle) DeleteById(c *gin.Context) {
	id := c.Param("id")
	/*
		query := "DELETE FROM footballstore WHERE id = $1"
		_, err := h.conn.Exec(context.Background(), query, id)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusNotFound, gin.H{"Error": "id not found"})
			return
		}
	*/
	c.JSON(http.StatusOK, gin.H{"message": "The good has been deleated"})
}
