package usecase

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type footballstore struct {
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

func (g *GoodsUsecase) ListStore(c *gin.Context) error {
	rows, err := g.conn.Query(context.Background(), `SELECT "id", "category", "name", "price" from "footballstore"`) //выполнение SQL запроса через соединение с бд
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when executing sql query"}) //ошибка при выполнении sql запроса
		return nil
	}
	defer rows.Close()

	var goods []footballstore //создается срез

	for rows.Next() { //метод, который используется для перебора строк
		var good footballstore                                              //объявляем переменную good, чтобы хранить данные для каждой строки, которую мы считваем из бд
		err := rows.Scan(&good.ID, &good.Category, &good.Name, &good.Price) //извлекаем данные из текущей струтуры базы данны(id, catogory ...) и присваиваем их полям структуры(&id, &category ...), чтобы можно было работать с этими данными
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when reading data"}) //Ошибка при чтении данных
			return nil
		}
		goods = append(goods, good) //добавляем извлеченные товары в наш срез
	}
	return nil
}

func (g *GoodsUsecase) UpdateStore(c *gin.Context, id string, updatedGoods *footballstore) error {
	query := "UPDATE footballstore SET category = $1, name = $2, price = $3 WHERE id = $4"                               //обновляем данные через SQL
	_, err := g.conn.Exec(context.Background(), query, updatedGoods.Category, updatedGoods.Name, updatedGoods.Price, id) //выполняем SQL запрос для обновления данных и присваиваем новые значения переменным
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error during the update"}) //ошибка при обнолении
		return err
	}
	return nil
}

func (g *GoodsUsecase) GetGoodByID(id string, good *footballstore) error {
	query := "SELECT id, name, category, price FROM footballstore WHERE id = $1" // выполняем SQL запрос, где выдается конкретный id, в данном случае 1
	row := g.conn.QueryRow(context.Background(), query, id)                      //используется чтобы выдать только 1 строку, в данном случае  id строку

	err := row.Scan(&good.ID, &good.Category, &good.Name, &good.Price) //Метод Scan извлекает значения из результата запроса и присваивает их полям структуры good.
	return err
}

func (g *GoodsUsecase) InsertStore(c *gin.Context, newGood *footballstore) error {
	query := "INSERT INTO footballstore (category, name, price) VALUES ($1, $2,  $3)"
	_, err := g.conn.Exec(context.Background(), query, newGood.Category, newGood.Name, newGood.Price)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't assign data"}) //не удалось присвоить данные
		return err
	}
	return nil
}

func (g *GoodsUsecase) DeleteById(c *gin.Context, id string) error {
	query := "DELETE FROM footballstore WHERE id = $1"
	_, err := g.conn.Exec(context.Background(), query, id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "id not found"})
		return err
	}
	return nil
}
