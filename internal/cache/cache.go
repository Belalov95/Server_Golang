package cache

import (
	"context"
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/repository"
	"sync"
)

type CacheDecorator struct {
	goodsRepo repository.GoodsProvider        //репозиторий для загрузки данных
	goods     map[string]models.Footballstore //добавляем в мапу id которые хранятся в массиве в качестве значения мапы
	mu        sync.RWMutex
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
	goods, ok := c.goods[id]
	c.mu.RUnlock()
	return goods, ok
}

// Set для сохранения товара в кеше
func (c *CacheDecorator) Set(goods models.Footballstore) {
	c.mu.Lock()
	c.goods[goods.ID] = goods
	c.mu.Unlock()
}

// удаляет товар из кеша
func (c *CacheDecorator) Delete(id string) {
	c.mu.Lock()
	delete(c.goods, id)
	c.mu.Unlock()
}

func (c *CacheDecorator) ListStore(ctx context.Context) ([]models.Footballstore, error) {
	return c.goodsRepo.ListStore(ctx) //загрузка данных из репозитория
}

func (c *CacheDecorator) UpdateStore(ctx context.Context, updaupdatedGoods *models.Footballstore) error {
	if err := c.goodsRepo.UpdateStore(ctx, updaupdatedGoods); err != nil {
		return err
	}
	c.mu.Lock()
	c.Set(*updaupdatedGoods) //обновляем кеш
	c.mu.Unlock()
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
	c.mu.Lock()
	c.Set(*goods) //сохраняем в кеше
	c.mu.Unlock()
	return goods, nil
}

func (c *CacheDecorator) InsertStore(ctx context.Context, newGood *models.Footballstore) error {
	if err := c.goodsRepo.InsertStore(ctx, newGood); err != nil {
		return err
	}
	c.mu.Lock()
	c.Set(*newGood) //сохраняем в кеше
	c.mu.Unlock()
	return nil
}

func (c *CacheDecorator) DeleteByID(ctx context.Context, id string) error {
	if err := c.goodsRepo.DeleteByID(ctx, id); err != nil {
		return err
	}
	c.mu.Lock()
	c.Delete(id) //удаляем из кеша
	c.mu.Unlock()
	return nil
}
