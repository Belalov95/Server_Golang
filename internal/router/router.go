package router

import (
	"example/web-service-gin/internal/handler"

	"github.com/gin-gonic/gin"
)

func GetRouter(handle *handler.Handle) *gin.Engine { // ф-ция, которая создает и настраивает HTTP-маршрутизатор

	router := gin.Default()                        //создаем маршрутизатор, который помогает серверу направлять входящие запросы к нужной ф-ции
	router.GET("/goods", handle.ListStore)         //определяем маршрут "/goods" и указываем ф-цию
	router.GET("/goods/:id", handle.GetGoodByID)   //определение маршрута, чтобы искать по id
	router.POST("/goods", handle.InsertStore)      //определение маршрута для POST-запроса на адрес /goods
	router.PUT("/goods/", handle.UpdateStore)      //определение маршрута для обновления данных о товаре
	router.DELETE("/goods/:id", handle.DeleteById) //определение маршрута для удаления

	return router
}
