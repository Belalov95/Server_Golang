// обработчик HTTP-запросов для веб-сервиса, который взаимодействует с базой данных для управления товарами в магазине
package handler

import (
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/usecase"
	"strconv"

	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
	//присваиваем переменной dbGoods список товаров
	dbGoods, err := h.goodsUC.ListStore(c)
	//проверка на ошибку после перебора
	if err != nil {
		wrappedErr := errors.Wrap(err, "failed to list goods in Liststore")
		slog.Error("ListStore goods error", slog.Any("error", wrappedErr))
		//Ошибка после перербора товаров
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error after sorting through the items"})
		return
	}
	//Преобразовываем данные из модели базы данных в структуру Response, чтобы вывести клиенту именно те данные из
	//базы, которые нужны ему
	var goods []models.GoodResponse
	for _, dbGood := range dbGoods {
		goods = append(goods, models.GoodResponse{
			ID:       dbGood.ID,
			Category: dbGood.Category,
			Name:     dbGood.Name,
			Price:    dbGood.Price,
		})
	}
	//отправляем список товаров в формате JSON
	c.JSON(http.StatusOK, goods)
}

// создаем эту ф-цию для обновления данных в магазине (в базе данных)
func (h *Handle) UpdateStore(c *gin.Context) {
	//создаем переменную updatedGoods чтобы хранить в ней обновленные товары
	updatedGoods := models.UpdateGoodRequest{}
	//считываем с помощью BindJSON новые данные которые отправил клиент и передаем их переменной updatedGoods
	if err := c.BindJSON(&updatedGoods); err != nil {
		wrappedErr := errors.Wrap(err, "failed to bind JSON in UpdateStore")
		slog.Error("UpdateStore BindJSON error", slog.Any("error", wrappedErr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "error unconnecting data"})
		return
	}
	dbModel := models.Footballstore{
		ID:       updatedGoods.ID,
		Category: updatedGoods.Category,
		Name:     updatedGoods.Name,
		Price:    updatedGoods.Price,
	}
	//Вызываем UseCase для обновления данных в базе
	if err := h.goodsUC.UpdateStore(c, &dbModel); err != nil {
		wrappedErr := errors.Wrap(err, "failed to update store in UpdateStore")
		slog.Error("UpdateStore uc.Updatestore error", slog.Any("error", wrappedErr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "error unconnecting data"})
		return
	}
	//отправляем клиенту ответ об успешном обновлении
	c.JSON(http.StatusOK, gin.H{"message": "The good has been successfully updated"})
}

// создается для поиска товара по его id
func (h *Handle) GetGoodByID(c *gin.Context) {
	//создаем для поиска товара по id
	idStr := c.Param("id")

	//Преобразовываем строку в целое число
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		wrappedErr := errors.Wrap(err, "invalid ID in GetGoodByID")
		slog.Error("Invalid ID", slog.Any("error", wrappedErr))
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid ID"})
		return
	}
	good, err := h.goodsUC.GetGoodByID(c, id)
	if err != nil {
		wrappedErr := errors.Wrap(err, "failed to get good by ID in GetGoodByID")
		slog.Error("GetGoodByID error", slog.Int("id", id), slog.Any("error", wrappedErr))
		c.JSON(http.StatusNotFound, gin.H{"error": "Good not found"})
		return
	}
	if good == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Good not found"})
		return
	}
	resp := models.GoodResponse{
		ID:       good.ID,
		Category: good.Category,
		Name:     good.Name,
		Price:    good.Price,
	}
	c.JSON(http.StatusOK, resp)
}

// добавляем новую запись в бд
func (h *Handle) InsertStore(c *gin.Context) {
	newGood := models.CreateGoodRequest{}
	//считываем json данные и присваиваем их переменной newGood
	if err := c.BindJSON(&newGood); err != nil {
		wrappedErr := errors.Wrap(err, "failed to bind JSON in InsertStore")
		slog.Error("InsertStore BindJSON error", slog.Any("error", wrappedErr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "couldn't assign data"})
		return
	}
	dbModel := models.Footballstore{
		Category: newGood.Category,
		Name:     newGood.Name,
		Price:    newGood.Price,
	}
	if err := h.goodsUC.InsertStore(c, &dbModel); err != nil {
		wrappedErr := errors.Wrap(err, "failed to insert store in InsertStore")
		slog.Error("InsertStore error", slog.Any("error", wrappedErr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "couldn't assign data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "The good added successfully"})
}

// Удаляем товар
func (h *Handle) DeleteById(c *gin.Context) {
	id := c.Param("id")
	if err := h.goodsUC.DeleteById(c, id); err != nil {
		wrappedErr := errors.Wrap(err, "failed to delete good by ID in DeleteById")
		slog.Error("DeleteById error", slog.Any("error", wrappedErr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "couldn't assign data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "The good has been deleated"})
}
