package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

var validate *validator.Validate

func Validate() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// Request для входных данных клиента
type CreateGoodRequest struct {
	Category string          `json:"category" validate:"required,min = 1,max = 50"`
	Name     string          `json:"name" validate:"required,min = 1,max = 100"`
	Price    decimal.Decimal `json:"price" validate:"required,gt = 0"`
}
type UpdateGoodRequest struct {
	ID       string          `json:"id" validate:"required"`
	Category string          `json:"category" validate:"omitempty,min = 1,max = 50"`
	Name     string          `json:"name" validate:"omitempty,min = 1,max = 100"`
	Price    decimal.Decimal `json:"price" validate:"omitempty,gt = 0"`
}

// Database Model для работы с бд
type Footballstore struct {
	ID       string          `db:"id"`       // объявление номера: serial primary key
	Category string          `db:"category"` // категория товара: одежда, обувь, аксессуары и т.д.
	Name     string          `db:"name"`     // название товара: форма, бутсы, брелки и т.д.
	Price    decimal.Decimal `db:"price"`    // цена товара
}

// Response для ответа клиенту
type GoodResponse struct {
	ID       string          `json:"id"`
	Category string          `json:"category"`
	Name     string          `json:"name"`
	Price    decimal.Decimal `json:"price"`
}

// Она берёт данные из модели БД и подготавливает их для ответа клиенту (например, по HTTP) .
func ToResponse(dbGoods []Footballstore) []GoodResponse {
	goods := make([]GoodResponse, len(dbGoods))
	// Преобразовываем данные из модели базы данных в структуру Response, чтобы вывести клиенту именно те данные из
	// базы, которые нужны ему
	for i, dbGood := range dbGoods {
		goods[i] = GoodResponse(dbGood)
	}
	return goods
}

// преобразовывает входные данные из UpdateGoodRequest в формат, который понимает бд Footballstore
func UpdatedGoodsDTO(updatedGoods UpdateGoodRequest) Footballstore {
	// ретерним чтобы возвращалось значение и мы могли использовать эту функцию
	return Footballstore(updatedGoods)
}
