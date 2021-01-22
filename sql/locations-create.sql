CREATE TABLE IF NOT EXISTS Locations (
    id SERIAL PRIMARY KEY,
    latitude INTEGER NOT NULL,
    longitude INTEGER NOT NULL,
    driver_id INTEGER NOT NULL
);