CREATE TABLE IF NOT EXISTS routes (
    id SERIAL PRIMARY KEY,
    order_id INT,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS route_points (
    id SERIAL PRIMARY KEY,
    route_id INT NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    cargo_id INT,
    address TEXT NOT NULL,
    planned_time TIMESTAMP,
    actual_time TIMESTAMP,
    status TEXT NOT NULL DEFAULT 'pending'
);