package cache

import (
	"context"
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/repository"
)

type CacheDecorator struct {
	goodsRepo repository.GoodsProvider        //репозиторий для загрузки данных
	goods     map[string]models.Footballstore //добавляем в мапу id которые хранятся в массиве в качестве значения мапы

}

func New(goodsRepo repository.GoodsProvider) *CacheDecorator {
	return &CacheDecorator{
		goodsRepo: goodsRepo,
		goods:     make(map[string]models.Footballstore),
	}
}

// метод Get для получения товара из кеша
func (c *CacheDecorator) Get(id string) (models.Footballstore, bool) {
	goods, ok := c.goods[id]
	return goods, ok
}

// Set для сохранения товара в кеше
func (c *CacheDecorator) Set(goods models.Footballstore) {
	c.goods[goods.ID] = goods
}

// удаляет товар из кеша
func (c *CacheDecorator) Delete(id string) {
	delete(c.goods, id)
}

func (c *CacheDecorator) ListStore(ctx context.Context) ([]models.Footballstore, error) {
	if len(c.goods) == 0 { //проверка кеша
		goods, err := c.goodsRepo.ListStore(ctx) //загрузка данных из репозитория
		if err != nil {
			return nil, err
		}
		//сохранение данных в кеше
		for _, good := range goods {
			c.goods[good.ID] = good
		}
		//вовзрат данных если кеш пустой
		return goods, nil
	}
	//возврат данных если кеш не пустой
	goodList := make([]models.Footballstore, 0, len(c.goods))
	for _, good := range c.goods {
		goodList = append(goodList, good)
	}
	return goodList, nil
}

func (c *CacheDecorator) UpdateStore(ctx context.Context, updaupdatedGoods *models.Footballstore) error {
	err := c.goodsRepo.UpdateStore(ctx, updaupdatedGoods)
	if err != nil {
		return err
	}
	c.Set(*updaupdatedGoods) //обновляем кеш
	return nil
}

func (c *CacheDecorator) GetGoodByID(ctx context.Context, id string) (*models.Footballstore, error) {
	if goods, ok := c.Get(id); ok {
		return &goods, nil
	}
	goods, err := c.goodsRepo.GetGoodByID(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Set(*goods) //сохраняем в кеше
	return goods, nil
}

func (c *CacheDecorator) InsertStore(ctx context.Context, newGood *models.Footballstore) error {
	err := c.goodsRepo.InsertStore(ctx, newGood)
	if err != nil {
		return err
	}
	c.Set(*newGood) //сохраняем в кеше
	return nil
}

func (c *CacheDecorator) DeleteByID(ctx context.Context, id string) error {
	err := c.goodsRepo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}
	c.Delete(id) //удаляем из кеша
	return nil
}
