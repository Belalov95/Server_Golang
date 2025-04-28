package repository

import (
	"example/web-service-gin/internal/apperr"
	"example/web-service-gin/internal/models"

	"github.com/pkg/errors"

	"github.com/jackc/pgx/v5"
	"golang.org/x/net/context"
)

type GoodsRepo struct {
	conn *pgx.Conn
}

func NewGoodsRepo(conn *pgx.Conn) *GoodsRepo {
	return &GoodsRepo{conn: conn}
}
func (g *GoodsRepo) ListStore(ctx context.Context) ([]models.Footballstore, error) {
	//выполнение SQL запроса через соединение с бд
	rows, err := g.conn.Query(ctx, `SELECT id, category, name, price FROM footballstore`)
	if err != nil {
		wrappedErr := errors.Wrap(err, "failed to execute query in ListStore")
		return nil, wrappedErr
	}
	defer rows.Close()
	//создается срез
	goods := []models.Footballstore{}

	//метод, который используется для перебора строк
	for rows.Next() {
		//объявляем переменную good, чтобы хранить данные для каждой строки, которую мы считваем из бд
		good := models.Footballstore{}
		//извлекаем данные из текущей струтуры базы данны(id, catogory ...) и присваиваем их полям структуры(&id,
		//&category ...), чтобы можно было работать с этими данными
		err := rows.Scan(&good.ID, &good.Category, &good.Name, &good.Price)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row in ListStore")
		}
		//добавляем извлеченные товары в наш срез
		goods = append(goods, good)
	}
	// Проверяем, не возникло ли ошибок при переборе строк
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows iteration error in ListStore")
	}
	return goods, nil
}
func (g *GoodsRepo) UpdateStore(ctx context.Context, updatedGoods *models.Footballstore) error {
	//обновляем данные через SQL
	query := "UPDATE footballstore SET category = $1, name = $2, price = $3 WHERE id = $4"
	//выполняем SQL запрос для обновления данных и присваиваем новые значения переменным
	_, err := g.conn.Exec(ctx, query, updatedGoods.ID, updatedGoods.Category, updatedGoods.Name, updatedGoods.Price)
	if err != nil {
		return errors.Wrap(err, "failed to update goods in UpdateStore")
	}
	return nil
}
func (g *GoodsRepo) GetGoodByID(ctx context.Context, id int) (*models.Footballstore, error) {
	// выполняем SQL запрос, где выдается конкретный id, в данном случае 1
	query := "SELECT id, name, category, price FROM Footballstore WHERE id = $1"
	//используется чтобы выдать только 1 строку, в данном случае  id строку
	row := g.conn.QueryRow(ctx, query, id)

	good := models.Footballstore{}
	//Метод Scan извлекает значения из результата запроса и присваивает их полям структуры good
	err := row.Scan(&good.ID, &good.Category, &good.Name, &good.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			//если не найдено, то возвращает nil
			return nil, apperr.ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to scan row in GetGoodByID")
	}
	return &good, nil
}
func (g *GoodsRepo) InsertStore(ctx context.Context, newGood *models.Footballstore) error {
	query := "INSERT INTO footballstore (category, name, price) VALUES ($1, $2,  $3)"
	_, err := g.conn.Exec(ctx, query, newGood.Category, newGood.Name, newGood.Price)
	if err != nil {
		return errors.Wrap(err, "failed to insert new good in InsertStore")
	}
	return nil
}
func (g *GoodsRepo) DeleteByID(ctx context.Context, id string) error {
	query := "DELETE FROM footballstore WHERE id = $1"
	_, err := g.conn.Exec(context.Background(), query, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete good by ID in DeleteByID")
	}
	return nil
}
