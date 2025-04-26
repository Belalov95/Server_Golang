package models

import (
	"github.com/shopspring/decimal"
)

// Request для входных данных клиента
type CreateGoodRequest struct {
	Category string          `json:"category"`
	Name     string          `json:"name"`
	Price    decimal.Decimal `json:"price"`
}
type UpdateGoodRequest struct {
	ID       string          `json:"id"`
	Category string          `json:"category"`
	Name     string          `json:"name"`
	Price    decimal.Decimal `json:"price"`
}

// Database Model для работы с бд
type Footballstore struct {
	ID       string          `db:"id"`       //объявление номера // serial primary key
	Category string          `db:"category"` //категория товара: Одежда, обувь, аксессуары и тд.
	Name     string          `db:"name"`     //название товара:форма, бутсы, брелки и тд.
	Price    decimal.Decimal `db:"price"`    //цена товара
}

// Response для ответа клиенту
type GoodResponse struct {
	ID       string          `json:"id"`
	Category string          `json:"category"`
	Name     string          `json:"name"`
	Price    decimal.Decimal `json:"price"`
}
