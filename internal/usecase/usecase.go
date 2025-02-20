package usecase

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
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
	rows, err := g.conn.Query(context.Background(), `SELECT "id", "category", "name", "price" from "footballstore"`) //выполнение SQL запроса через соединение с бд
	if err != nil {
		slog.Error("ListStore Query error", slog.Any("error", err))
		return nil, err
	}
	defer rows.Close()

	goods := []Footballstore{} //создается срез
	for rows.Next() {          //метод, который используется для перебора строк
		good := Footballstore{}                                                              //объявляем переменную good, чтобы хранить данные для каждой строки, которую мы считваем из бд
		if err := rows.Scan(&good.ID, &good.Category, &good.Name, &good.Price); err != nil { //извлекаем данные из текущей струтуры базы данны(id, catogory ...) и присваиваем их полям структуры(&id, &category ...), чтобы можно было работать с этими данными
			slog.Error("ListStore Scan error", slog.Any("error", err))
			return nil, err
		}
		goods = append(goods, good) //добавляем извлеченные товары в наш срез
	}
	return goods, nil
}

func (g *GoodsUsecase) GetGoodByID(ctx context.Context, id string) (*Footballstore, error) {
	query := "SELECT id, name, category, price FROM footballstore WHERE id = $1" // выполняем SQL запрос, где выдается конкретный id, в данном случае 1
	row := g.conn.QueryRow(context.Background(), query, id)                      //используется чтобы выдать только 1 строку, в данном случае  id строку

	good := Footballstore{}
	err := row.Scan(&good.ID, &good.Category, &good.Name, &good.Price) //Метод Scan извлекает значения из результата запроса и присваивает их полям структуры good.
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil //если не найдено, то возвращается nil
		}
		slog.Error("GetGoodByID Scan error", slog.Any("error", err))
		return nil, err
	}
	return &good, nil
}

func (g *GoodsUsecase) DeleteById(ctx context.Context, id string) error {
	query := "DELETE FROM footballstore WHERE id = $1"
	_, err := g.conn.Exec(context.Background(), query, id)
	if err != nil {
		slog.Error("DeleteById Exec error", slog.Any("error", err))
		return err
	}
	return nil
}

// TODO: редачить

func (g *GoodsUsecase) UpdateStore(ctx context.Context, good *Footballstore) error {
	//обновляем данные через SQL
	//выполняем SQL запрос для обновления данных и присваиваем новые значения переменным
	query := "UPDATE footballstore SET category = $1, name = $2, price = $3 WHERE id = $4"
	if _, err := g.conn.Exec(ctx, query, good.Category, good.Name, good.Price, good.ID); err != nil {
		slog.Error("ListStore Query error", slog.Any("error", err))
		return err
	}
	return nil
}

func (g *GoodsUsecase) InsertStore(c *gin.Context, newGood *Footballstore) error {
	query := "INSERT INTO footballstore (category, name, price) VALUES ($1, $2,  $3)"
	_, err := g.conn.Exec(context.Background(), query, newGood.Category, newGood.Name, newGood.Price)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't assign data"}) //не удалось присвоить данные
		return err
	}
	return nil
}
