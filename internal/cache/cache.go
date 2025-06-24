package cache

import (
	"context"
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/repository"
	"sync"
)

type CacheDecorator struct {
	goodsRepo repository.GoodsProvider //репозиторий для загрузки данных
	mu        sync.RWMutex
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
	c.mu.RLock()
	defer c.mu.RUnlock()
	goods, ok := c.goods[id]
	return goods, ok
}

// Set для сохранения товара в кеше
func (c *CacheDecorator) Set(goods models.Footballstore) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.goods[goods.ID] = goods
}

// удаляет товар из кеша
func (c *CacheDecorator) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.goods, id)
}

func (c *CacheDecorator) ListStore(ctx context.Context) ([]models.Footballstore, error) {
	return c.goodsRepo.ListStore(ctx) //загрузка данных из репозитория
}

func (c *CacheDecorator) UpdateStore(ctx context.Context, updaupdatedGoods *models.Footballstore) error {
	if err := c.goodsRepo.UpdateStore(ctx, updaupdatedGoods); err != nil {
		return err
	}
	c.Set(*updaupdatedGoods) //обновляем кеш
	return nil
}

func (c *CacheDecorator) GetGoodByID(ctx context.Context, id string) (*models.Footballstore, error) {
	if goods, ok := c.Get(id); ok { //проверяем содержится ли товар в кеше
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
	if err := c.goodsRepo.InsertStore(ctx, newGood); err != nil {
		return err
	}
	c.Set(*newGood) //сохраняем в кеше
	return nil
}

func (c *CacheDecorator) DeleteByID(ctx context.Context, id string) error {
	if err := c.goodsRepo.DeleteByID(ctx, id); err != nil {
		return err
	}
	c.Delete(id) //удаляем из кеша
	return nil
}
