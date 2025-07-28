package usecase

import (
	"context"
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/repository"
)

type GoodsUsecase struct {
	goodsRepo repository.GoodsProvider
}

func New(goodsRepo repository.GoodsProvider) *GoodsUsecase {
	return &GoodsUsecase{goodsRepo: goodsRepo}
}

func (g *GoodsUsecase) ListStore(ctx context.Context) ([]models.Footballstore, error) {
	return g.goodsRepo.ListStore(ctx)
}

func (g *GoodsUsecase) UpdateStore(ctx context.Context, updatedGoods *models.Footballstore) error {
	return g.goodsRepo.UpdateStore(ctx, updatedGoods)
}

func (g *GoodsUsecase) GetGoodByID(ctx context.Context, id string) (*models.Footballstore, error) {
	return g.goodsRepo.GetGoodByID(ctx, id)
}

func (g *GoodsUsecase) InsertStore(ctx context.Context, newGood *models.Footballstore) error {
	return g.goodsRepo.InsertStore(ctx, newGood)
}

func (g *GoodsUsecase) DeleteById(ctx context.Context, id string) error {
	return g.goodsRepo.DeleteByID(ctx, id)
}
