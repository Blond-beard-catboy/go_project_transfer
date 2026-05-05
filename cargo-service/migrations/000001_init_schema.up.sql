CREATE TABLE IF NOT EXISTS cargos (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    weight DOUBLE PRECISION NOT NULL,
    pickup_location TEXT NOT NULL,
    delivery_location TEXT NOT NULL,
    pickup_date TIMESTAMP NOT NULL,
    delivery_date TIMESTAMP NOT NULL,
    owner_id INT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);