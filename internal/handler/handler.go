// обработчик HTTP-запросов для веб-сервиса, который взаимодействует с базой данных для управления товарами в магазине
package handler

import (
	"example/web-service-gin/internal/apperr"
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/usecase"

	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		slog.Error("ListStore goods error", slog.Any("error", err))
		//Ошибка после перербора товаров
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error after sorting through the items"})
		return
	}
	resp := models.ToResponse(dbGoods) //goods содержит результат функции toresponse
	//отправляем список товаров в формате JSON
	c.JSON(http.StatusOK, resp)
	return
}

// создаем эту ф-цию для обновления данных в магазине (в базе данных)
func (h *Handle) UpdateStore(c *gin.Context) {
	id := c.Param("id")
	//создаем переменную updatedGoods чтобы хранить в ней обновленные товары
	updatedGoods := models.UpdateGoodRequest{}
	//считываем с помощью BindJSON новые данные которые отправил клиент и передаем их переменной updatedGoods
	if err := c.BindJSON(&updatedGoods); err != nil {
		slog.Error("UpdateStore BindJSON error", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error unconnecting data"})
		return
	}

	//устанавливаем ID из URL в структуру обновления
	updatedGoods.ID = id
	// валидация
	if err := models.ValidateStruct(updatedGoods); err != nil {
		errors := make(map[string]string)
		for _, fieldErr := range err.(validator.ValidationErrors) {
			errors[fieldErr.Field()] = fieldErr.Tag() //записываем в мапу где fied - name, price ... и tag - required, gt ...
		}
		c.JSON(http.StatusBadRequest, gin.H{"validation error": errors})
		return
	}

	//обновляем dbModel с помощью updatedGoodsTDO и объявляем переменную чтобы в дальнейшем использовать ее
	dbModel := models.UpdatedGoodsDTO(updatedGoods)

	//Вызываем UseCase для обновления данных в базе
	if err := h.goodsUC.UpdateStore(c, &dbModel); err != nil {
		slog.Error("UpdateStore uc.Updatestore error", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "error unconnecting data"})
		return
	}
	//отправляем клиенту ответ об успешном обновлении
	c.JSON(http.StatusOK, gin.H{"message": "The good has been successfully updated"})
	return
}

// создается для поиска товара по его id
func (h *Handle) GetGoodByID(c *gin.Context) {
	//создаем для поиска товара по id
	id := c.Param("id")

	good, err := h.goodsUC.GetGoodByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found"})
			return
		}
		//логируем ошибку
		slog.Error("GetGoodByID error", slog.String("id", id), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "500 Internal Server Error"})
		return
	}

	//преобразовываем товар в формат ответа
	resp := models.GoodResponse{
		ID:       good.ID,
		Category: good.Category,
		Name:     good.Name,
		Price:    good.Price,
	}
	c.JSON(http.StatusOK, resp)
	return
}

// добавляем новую запись в бд
func (h *Handle) InsertStore(c *gin.Context) {
	newGood := models.CreateGoodRequest{}
	//считываем json данные и присваиваем их переменной newGood
	if err := c.BindJSON(&newGood); err != nil {
		slog.Error("InsertStore BindJSON error", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "couldn't assign data"})
		return
	}

	//валидация
	if err := models.ValidateStruct(newGood); err != nil {
		errors := make(map[string]string)
		for _, fieldErr := range err.(validator.ValidationErrors) {
			errors[fieldErr.Field()] = fieldErr.Tag()
		}
		c.JSON(http.StatusBadRequest, gin.H{"validation error": errors})
		return
	}

	dbModel := models.Footballstore{
		Category: newGood.Category,
		Name:     newGood.Name,
		Price:    newGood.Price,
	}
	if err := h.goodsUC.InsertStore(c, &dbModel); err != nil {
		slog.Error("InsertStore error", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "couldn't assign data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "The good added successfully"})
	return
}

// Удаляем товар
func (h *Handle) DeleteById(c *gin.Context) {
	id := c.Param("id")
	if err := h.goodsUC.DeleteById(c, id); err != nil {
		slog.Error("DeleteById error", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "couldn't assign data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "The good has been deleated"})
	return
}
