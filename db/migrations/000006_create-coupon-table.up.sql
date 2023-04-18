CREATE TABLE IF NOT EXISTS coupon (
    id SERIAL PRIMARY KEY,

    discount REAL NOT NULL,
    description TEXT NOT NULL,

    request_product_id INTEGER NOT NULL REFERENCES request_product(id)
);
