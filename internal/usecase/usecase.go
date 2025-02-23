package usecase

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type Footballstore struct {
	ID       string          `json:"id"`       //объявление номера // serial primary key
	Category string          `json:"category"` //категория товара: Одежда, обувь, аксессуары и тд.
	Name     string          `json:"name"`     //название товара:форма, бутсы, брелки и тд.
	Price    decimal.Decimal `json:"price"`    //цена товара
}

type GoodsUsecase struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *GoodsUsecase {
	return &GoodsUsecase{conn: conn}
}

func (g *GoodsUsecase) ListStore(ctx context.Context) ([]Footballstore, error) {
	//выполнение SQL запроса через соединение с бд
	rows, err := g.conn.Query(ctx, `SELECT "id", "category", "name", "price" from "footballstore"`)
	if err != nil {
		slog.Error("ListStore Query error", slog.Any("error", err))
		return nil, err
	}
	defer rows.Close()
	//создается срез
	goods := []Footballstore{}

	//метод, который используется для перебора строк
	for rows.Next() {
		//объявляем переменную good, чтобы хранить данные для каждой строки, которую мы считваем из бд
		good := Footballstore{}
		//извлекаем данные из текущей струтуры базы данны(id, catogory ...) и присваиваем их полям структуры(&id,
		//&category ...), чтобы можно было работать с этими данными
		err := rows.Scan(&good.ID, &good.Category, &good.Name, &good.Price)
		if err != nil {
			slog.Error("ListStore Scan error", slog.Any("error", err))
			return nil, err
		}
		//добавляем извлеченные товары в наш срез
		goods = append(goods, good)
	}
	return goods, nil
}

func (g *GoodsUsecase) UpdateStore(ctx context.Context, good *Footballstore) error {
	//обновляем данные через SQL
	query := "UPDATE footballstore SET category = $1, name = $2, price = $3 WHERE id = $4"
	//выполняем SQL запрос для обновления данных и присваиваем новые значения переменным
	_, err := g.conn.Exec(ctx, query, good.Category, good.Name, good.Price)
	if err != nil {
		slog.Error(" UpdateStore goods error", slog.Any("error", err))
		return err
	}
	return nil
}

func (g *GoodsUsecase) GetGoodByID(ctx context.Context, id string) (*Footballstore, error) {
	// выполняем SQL запрос, где выдается конкретный id, в данном случае 1
	query := "SELECT id, name, category, price FROM footballstore WHERE id = $1"
	//используется чтобы выдать только 1 строку, в данном случае  id строку
	row := g.conn.QueryRow(context.Background(), query, id)

	good := Footballstore{}
	//Метод Scan извлекает значения из результата запроса и присваивает их полям структуры good
	err := row.Scan(&good.ID, &good.Category, &good.Name, &good.Price)
	if err != nil {
		if err == pgx.ErrNoRows {
			//если не найдено, то возвращает nil
			return nil, nil
		}
		slog.Error("GetGoodByID Scan error", slog.Any("error", err))
		return nil, err
	}
	return &good, nil
}
func (g *GoodsUsecase) InsertStore(ctx context.Context, newGood *Footballstore) error {
	query := "INSERT INTO footballstore (category, name, price) VALUES ($1, $2,  $3)"
	_, err := g.conn.Exec(ctx, query, newGood.Category, newGood.Name, newGood.Price)
	if err != nil {
		slog.Error("InsertStore error", slog.Any("error", err))
		return err
	}
	return nil
}

func (g *GoodsUsecase) DeleteById(ctx context.Context, id string) error {
	query := "DELETE FROM footballstore WHERE id = $1"
	_, err := g.conn.Exec(context.Background(), query, id)
	if err != nil {
		slog.Error("DeleteByID Scan error", slog.Any("error", err))
		return err
	}
	return nil
}
