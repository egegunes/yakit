DROP TABLE IF EXISTS brands;

CREATE TABLE brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256)
);

DROP TABLE IF EXISTS models;

CREATE TABLE models (
    id SERIAL PRIMARY KEY,
    brand_id INTEGER NOT NULL,
    name VARCHAR(256)
);

DROP TABLE IF EXISTS vehicles;

CREATE TABLE vehicles (
    id SERIAL PRIMARY KEY,
    model_id INTEGER NOT NULL,
    year INTEGER NOT NULL
);

DROP TABLE IF EXISTS entries;

CREATE TABLE entries (
    id SERIAL PRIMARY KEY,
    vehicle_id INTEGER NOT NULL,
    consumption INTEGER NOT NULL,
    message VARCHAR(256)
);

