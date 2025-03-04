// обработчик HTTP-запросов для веб-сервиса, который взаимодействует с базой данных для управления товарами в магазине
package handler

import (
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/usecase"

	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handle struct {
	// goodsUC это поле структуры Handle и это просто название. А usecase.GoodsProvider это тип поля goodsUC
	//В данном случае это интерфейс, определенный в пакете usecase
	goodsUC usecase.GoodsProvider
}

func New(goodsUC usecase.GoodsProvider) *Handle {
	return &Handle{goodsUC: goodsUC}
}

// создается для получения данных из таблицы бд
func (h *Handle) ListStore(c *gin.Context) {
	//присваиваем переменной goods список товаров
	goods, err := h.goodsUC.ListStore(c)
	//проверка на ошибку после перебора
	if err != nil {
		slog.Error("ListStore goods error", slog.Any("error", err))
		//Ошибка после перербора товаров
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error after sorting through the items"})
		return
	}
	//отправляем список товаров в формате JSON
	c.JSON(http.StatusOK, goods)
}

// создаем эту ф-цию для обновления данных в магазине (в базе данных)
func (h *Handle) UpdateStore(c *gin.Context) {
	//создаем переменную updatedGoods чтобы хранить в ней обновленные товары
	updatedGoods := models.Footballstore{}
	//считываем с помощью BindJSON новые данные которые отправил клиент и передаем их переменной updatedGoods
	if err := c.BindJSON(&updatedGoods); err != nil {
		slog.Error("UpdateStore BindJSON error", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "error unconnecting data"})
		return
	}
	if err := h.goodsUC.UpdateStore(c, &updatedGoods); err != nil {
		slog.Error("UpdateStore uc.Updatestore error", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "error unconnecting data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "The good has been successfully updated"}) //отправляем клиенту ответ об успешном обновлении
}

// создается для поиска товара по его id
func (h *Handle) GetGoodByID(c *gin.Context) {
	//создаем для поиска товара по id
	id := c.Param("id")
	good, err := h.goodsUC.GetGoodByID(c, id)
	if err != nil {
		slog.Error("GetGoodByID error", slog.Any("error", err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Good not found"})
		return
	}
	if good != nil {
		slog.Error("Good not found", slog.Any("error", err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Good not found"})
		return
	}
	c.JSON(http.StatusOK, good)
}

// добавляем новую запись в бд
func (h *Handle) InsertStore(c *gin.Context) {
	newGood := models.Footballstore{}
	//считываем json данные и присваиваем их переменной newGood
	if err := c.BindJSON(&newGood); err != nil {
		slog.Error("InsertStore BindJSON error", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "couldn't assign data"})
		return
	}
	if err := h.goodsUC.InsertStore(c, &newGood); err != nil {
		slog.Error("InsertStore error", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "couldn't assign data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "The good added successfully"})
}
func (h *Handle) DeleteById(c *gin.Context) {
	id := c.Param("id")
	if err := h.goodsUC.DeleteById(c, id); err != nil {
		slog.Error("DeleteById error", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "couldn't assign data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "The good has been deleated"})
}
