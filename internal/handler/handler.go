package handler

import (
	"example/web-service-gin/internal/usecase"
	"log"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handle struct {
	goodsUC *usecase.GoodsUsecase
}

func New(goodsUC *usecase.GoodsUsecase) *Handle {
	return &Handle{goodsUC: goodsUC}
}

// создается для получения данных из таблицы бд
func (h *Handle) ListStore(c *gin.Context) {
	goods, err := h.goodsUC.ListStore(c) //присваиваем переменной goods список товаров
	if err != nil {                      //проверка на ошибку после перебора
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error after sorting through the items"}) //Ошибка после перербора товаров
		return
	}
	c.JSON(http.StatusOK, goods) //отправляем список товаров в формате JSON
}

// Создаем эту ф-цию для обновления данных в магазине (в базе данных)
//
//	PUT `http://localhost/goods` тело запроса: {
//	   "id": ""
//	   "category": ""
//	   "name": ""
//	   "price": ""
//	}
func (h *Handle) UpdateStore(c *gin.Context) {
	//создаем переменную updatedGoods чтобы хранить в ней обновленные товары
	var updatedGoods usecase.Footballstore
	//считываем с помощью BindJSON новые данные которые отправил клиент и передаем их переменной updatedGoods
	if err := c.BindJSON(&updatedGoods); err != nil {
		slog.Error("UpdateStore BindJSON error", slog.Any("error", err))
		//некоректные данные
		c.JSON(http.StatusBadRequest, gin.H{"error": "error unconnecting data"})
		return
	}

	if err := h.goodsUC.UpdateStore(c, &updatedGoods); err != nil {
		slog.Error("UpdateStore uc.UpdateStore error", slog.Any("error", err))
		//некоректные данные
		c.JSON(http.StatusBadRequest, gin.H{"error": "error goods updating"})
		return
	}

	//отправляем клиенту ответ об успешном обновлении
	c.JSON(http.StatusOK, gin.H{"message": "the good has been successfully updated"})
}

// добавляем новую запись в бд
func (h *Handle) InsertStore(c *gin.Context) {
	var newGood usecase.Footballstore

	if err := c.BindJSON(&newGood); err != nil { //считываем json данные и присваиваем их переменной newGood
		slog.Error("InsertStore BindJSON error", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "couldn't assign data"})
		return
	}

	if err := h.goodsUC.InsertStore(c, &newGood); err != nil {
		slog.Error("InsertStore error", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't add the good"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "The good added successfully"})
}

// создается для поиска товара по его id
func (h *Handle) GetGoodByID(c *gin.Context) {
	//создаем для поиска товара по id
	id := c.Param("id")
	good, err := h.goodsUC.GetGoodByID(c, id)
	if err != nil {
		slog.Error("GetGoodByID error", slog.Any("error", err))
		//товар не найден
		c.JSON(http.StatusNotFound, gin.H{"error": "Good not found"})
		return
	}
	c.JSON(http.StatusOK, good)
}

func (h *Handle) DeleteById(c *gin.Context) {
	id := c.Param("id")
	if err := h.goodsUC.DeleteById(c, id); err != nil {
		slog.Error("DeleteById error", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't delete the good"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "The good has been deleated"})
}
