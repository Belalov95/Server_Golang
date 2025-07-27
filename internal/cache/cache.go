package cache

import (
	"context"
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/repository"
	"sync"
	"time"
)

// обертка для хранения товара и времени последнего обновления
type wrapGoods struct {
	good      models.Footballstore
	updatedAt time.Time
}

type CacheDecorator struct {
	goodsRepo repository.GoodsProvider // репозиторий для загрузки данных
	mu        sync.RWMutex
	goods     map[string]wrapGoods // добавляем в мапу id которые хранятся в массиве в качестве значения мапы
	ttl       time.Duration        // время жизни кеша
}

func New(goodsRepo repository.GoodsProvider, ttl time.Duration) *CacheDecorator {
	return &CacheDecorator{
		goodsRepo: goodsRepo,
		goods:     make(map[string]wrapGoods),
		ttl:       ttl,
	}
}

// метод Get для получения товара из кеша (чтение)
func (c *CacheDecorator) Get(id string) (models.Footballstore, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	wrapped, ok := c.goods[id]
	if !ok {
		return models.Footballstore{}, false
	}

	// проверка времени жизни записи
	if time.Since(wrapped.updatedAt) > c.ttl {
		c.mu.RUnlock()
		c.mu.Lock()
		delete(c.goods, id)
		c.mu.Unlock()
		c.mu.RLock()
		return models.Footballstore{}, false
	}
	return wrapped.good, true
}

// Set для сохранения товара в кеше
func (c *CacheDecorator) Set(goods models.Footballstore) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.goods[goods.ID] = wrapGoods{
		good:      goods,
		updatedAt: time.Now(),
	}
}

// удаляет товар из кеша
func (c *CacheDecorator) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.goods, id)
}

// удаляем просроченные записи из кеша
func (c *CacheDecorator) CleanupExpired() {
	now := time.Now() // фиксируем время старта
	c.mu.Lock()
	defer c.mu.Unlock()
	for id, wrapped := range c.goods {
		if now.Sub(wrapped.updatedAt) > c.ttl {
			delete(c.goods, id)
		}
	}
}

// запускаем горутину для регулярной очистки кеша
func (c *CacheDecorator) RunCleanup(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select { // ждет пока один из каналов не даст сигнал и запускает соответствующий код
			case <-ticker.C: // канал на который тикер каждый интервал времени отправляет сигнал
				c.CleanupExpired() // если пришел сигнал от тикера то выполняем очистку
			case <-ctx.Done(): // сигнал, который отправляется, когда контекст отменяется
				return
			}
		}
	}()
}

func (c *CacheDecorator) ListStore(ctx context.Context) ([]models.Footballstore, error) {
	return c.goodsRepo.ListStore(ctx) // загрузка данных из репозитория
}

func (c *CacheDecorator) UpdateStore(ctx context.Context, updaupdatedGoods *models.Footballstore) error {
	if err := c.goodsRepo.UpdateStore(ctx, updaupdatedGoods); err != nil {
		return err
	}
	c.Set(*updaupdatedGoods) // обновляем кеш
	return nil
}

func (c *CacheDecorator) GetGoodByID(ctx context.Context, id string) (*models.Footballstore, error) {
	if goods, ok := c.Get(id); ok { // проверяем содержится ли товар в кеше
		return &goods, nil
	}
	goods, err := c.goodsRepo.GetGoodByID(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Set(*goods) // сохраняем в кеше
	return goods, nil
}

func (c *CacheDecorator) InsertStore(ctx context.Context, newGood *models.Footballstore) error {
	if err := c.goodsRepo.InsertStore(ctx, newGood); err != nil {
		return err
	}
	c.Set(*newGood) // сохраняем в кеше
	return nil
}

func (c *CacheDecorator) DeleteByID(ctx context.Context, id string) error {
	if err := c.goodsRepo.DeleteByID(ctx, id); err != nil {
		return err
	}
	c.Delete(id) // удаляем из кеша
	return nil
}
