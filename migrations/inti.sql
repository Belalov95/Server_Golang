CREATE TABLE Footballstore (
    id       SERIAL PRIMARY KEY,          -- Автоматически увеличиваемый уникальный идентификатор
    category VARCHAR(255) NOT NULL,       -- Категория товара (строка с ограничением длины)
    name     VARCHAR(255) NOT NULL,       -- Название товара (строка с ограничением длины)
    price    NUMERIC(10, 2) NOT NULL      -- Цена товара (десятичная точность: до 10 цифр, 2 после запятой)
);