package usecase

import (
	"context"
	"example/web-service-gin/internal/models"
)

type GoodsProvider interface {
	ListStore(ctx context.Context) ([]models.Footballstore, error)
	UpdateStore(ctx context.Context, updatedGoods *models.Footballstore) error
	GetGoodByID(ctx context.Context, id string) (*models.Footballstore, error)
	InsertStore(ctx context.Context, newGood *models.Footballstore) error
	DeleteById(ctx context.Context, id string) error
}
