CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    cargo_id INT NOT NULL,
    customer_id INT NOT NULL,
    driver_id INT,
    route_id INT NOT NULL,
    status TEXT NOT NULL,
    contract_file TEXT,
    created_at TIMESTAMP NOT NULL
);