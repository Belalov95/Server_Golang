-- +goose Up
CREATE TABLE IF NOT EXISTS footballstore (
    id SERIAL PRIMARY KEY,
    category VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    price NUMERIC(10, 2) NOT NULL
);

INSERT INTO footballstore (category, name, price) VALUES
  ('Одежда', 'Футболка сборной Аргентины', 3499.00),
  ('Обувь', 'Бутсы Nike Mercurial', 7999.99),
  ('Аксессуары', 'Бутылка для воды Adidas', 999.50),
  ('Одежда', 'Шорты тренировочные', 1999.00),
  ('Аксессуары', 'Брелок в форме мяча', 299.00);

-- +goose Down
DROP TABLE IF EXISTS footballstore;