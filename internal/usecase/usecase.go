package usecase

import (
	"context"
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/repository"

	"github.com/jackc/pgx/v5"
)

type GoodsUsecase struct {
	goodsRepo *repository.GoodsRepo
}

func New(conn *pgx.Conn) *GoodsUsecase {
	return &GoodsUsecase{goodsRepo: repository.NewGoodsRepo(conn)}
}

func (g *GoodsUsecase) ListStore(ctx context.Context) ([]models.Footballstore, error) {
	return g.goodsRepo.ListStore(ctx)
}

func (g *GoodsUsecase) UpdateStore(ctx context.Context, updatedGoods *models.Footballstore) error {
	return g.goodsRepo.UpdateStore(ctx, updatedGoods)
}

func (g *GoodsUsecase) GetGoodByID(ctx context.Context, id int) (*models.Footballstore, error) {
	return g.goodsRepo.GetGoodByID(ctx, id)
}
func (g *GoodsUsecase) InsertStore(ctx context.Context, newGood *models.Footballstore) error {
	return g.goodsRepo.InsertStore(ctx, newGood)
}

func (g *GoodsUsecase) DeleteById(ctx context.Context, id string) error {
	return g.goodsRepo.DeleteByID(ctx, id)
}
