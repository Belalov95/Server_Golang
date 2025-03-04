package models

import "github.com/shopspring/decimal"

type Footballstore struct {
	ID       string          `json:"id"`       //объявление номера // serial primary key
	Category string          `json:"category"` //категория товара: Одежда, обувь, аксессуары и тд.
	Name     string          `json:"name"`     //название товара:форма, бутсы, брелки и тд.
	Price    decimal.Decimal `json:"price"`    //цена товара
}
