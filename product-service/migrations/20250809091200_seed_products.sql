-- +goose Up
INSERT INTO products (name, description, price, stock)
VALUES
  ('Keyboard', 'Mechanical keyboard', 150.00, 50),
  ('Mouse', 'Wireless mouse', 45.50, 200),
  ('Monitor', '27 inch 144Hz', 320.00, 20);
-- +goose Down
TRUNCATE products;
