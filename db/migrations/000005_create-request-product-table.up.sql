CREATE TABLE IF NOT EXISTS request_product (
    id SERIAL PRIMARY KEY,
    product_value REAL NOT NULL,

    request_id INTEGER NOT NULL REFERENCES request(id),
    product_id VARCHAR NOT NULL REFERENCES product(id)
);
