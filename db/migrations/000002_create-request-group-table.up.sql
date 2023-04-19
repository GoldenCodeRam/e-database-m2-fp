CREATE TABLE IF NOT EXISTS status_type (
    status_type VARCHAR PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS request (
    id SERIAL PRIMARY KEY,

    person_id INTEGER NOT NULL REFERENCES person(id),
    status_type VARCHAR NOT NULL REFERENCES status_type(status_type),
    date TIMESTAMP NOT NULL,
    country VARCHAR NOT NULL
);
