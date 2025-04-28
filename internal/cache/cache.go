package cache

import (
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/repository"
)

type CacheDecorator struct {
	goodsRepo repository.GoodsProvider

	goods map[id]models.Footballstore
}

// метод Get для получения товара из кеша
func (c *CacheDecorator) Get(id int) (goods, bool) {
	goods, ok := c.goods[id]
	return goods, ok
}

// Set для сохранения товара в кеше
func (c *CacheDecorator) Set(goods) {

}
