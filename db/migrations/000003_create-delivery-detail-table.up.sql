CREATE TABLE IF NOT EXISTS delivery_detail (
    id SERIAL PRIMARY KEY,
    address VARCHAR NOT NULL,

    description TEXT,

    request_id INTEGER NOT NULL REFERENCES request(id)
);
