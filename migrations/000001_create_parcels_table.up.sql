CREATE TABLE parcels (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL,
    pickup_address TEXT NOT NULL,
    dropoff_address TEXT NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    status VARCHAR(255) DEFAULT 'Created',
    pickup_time TIMESTAMP,
    delivery_time TIMESTAMP,
    motorbike_id INT
);