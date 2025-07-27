package cache_test

import (
	"context"
	"testing"
	"time"

	"example/web-service-gin/internal/cache"
	"example/web-service-gin/internal/mocks"
	"example/web-service-gin/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCacheDecorator_GetGoodByID(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name         string
		cacheTTL     time.Duration
		cachedBefore bool
		wantFromRepo bool
	}
	testGood := &models.Footballstore{
		ID:    "123",
		Name:  "Ball",
		Price: decimal.NewFromInt(1000),
	}

	tests := []testCase{
		{
			name:         "получаем из репозитория, кеш пуст",
			cacheTTL:     time.Minute,
			cachedBefore: false,
			wantFromRepo: true,
		},
		{
			name:         "получаем из кеша, он актуален",
			cacheTTL:     time.Minute,
			cachedBefore: true,
			wantFromRepo: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockGoodsProvider(ctrl)

			if tc.wantFromRepo {
				mockRepo.EXPECT().
					GetGoodByID(gomock.Any(), testGood.ID).
					Return(testGood, nil).
					Times(1)
			}

			c := cache.New(mockRepo, tc.cacheTTL)
			if tc.cachedBefore {
				c.Set(*testGood)
			}

			got, err := c.GetGoodByID(context.Background(), testGood.ID)
			assert.NoError(t, err)
			assert.Equal(t, testGood, got)
		})
	}
}

func TestCacheDecorator_ListStore(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockGoodsProvider(ctrl)
	expected := []models.Footballstore{{ID: "1", Name: "Boots", Price: decimal.NewFromInt(2000)}}

	mockRepo.EXPECT().
		ListStore(gomock.Any()).
		Return(expected, nil).
		Times(1)

	c := cache.New(mockRepo, time.Minute)
	got, err := c.ListStore(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestCacheDecorator_InsertStore(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockGoodsProvider(ctrl)
	good := &models.Footballstore{ID: "2", Name: "Jersey", Price: decimal.NewFromInt(1500)}

	mockRepo.EXPECT().
		InsertStore(gomock.Any(), good).
		Return(nil).
		Times(1)

	c := cache.New(mockRepo, time.Minute)
	err := c.InsertStore(context.Background(), good)
	assert.NoError(t, err)

	cached, ok := c.Get(good.ID)
	assert.True(t, ok)
	assert.Equal(t, *good, cached)
}

func TestCacheDecorator_UpdateStore(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockGoodsProvider(ctrl)
	good := &models.Footballstore{ID: "3", Name: "Gloves", Price: decimal.NewFromInt(800)}

	mockRepo.EXPECT().
		UpdateStore(gomock.Any(), good).
		Return(nil).
		Times(1)

	c := cache.New(mockRepo, time.Minute)
	err := c.UpdateStore(context.Background(), good)
	assert.NoError(t, err)

	cached, ok := c.Get(good.ID)
	assert.True(t, ok)
	assert.Equal(t, *good, cached)
}

func TestCacheDecorator_DeleteByID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockGoodsProvider(ctrl)
	id := "4"
	good := models.Footballstore{ID: id, Name: "Cap", Price: decimal.NewFromInt(300)}

	mockRepo.EXPECT().
		DeleteByID(gomock.Any(), id).
		Return(nil).
		Times(1)

	c := cache.New(mockRepo, time.Minute)
	c.Set(good) // кэшируем заранее

	err := c.DeleteByID(context.Background(), id)
	assert.NoError(t, err)

	_, ok := c.Get(id)
	assert.False(t, ok) // должно быть удалено
}
