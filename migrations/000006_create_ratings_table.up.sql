CREATE TABLE ratings (
                         id SERIAL PRIMARY KEY,
                         sender_id INT NOT NULL,
                         motorbike_id INT NOT NULL,
                         parcel_id INT NOT NULL,
                         rating INT CHECK (rating >= 1 AND rating <= 5), -- Rating must be between 1 and 5
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         UNIQUE (parcel_id), -- Ensure one rating per parcel
                         FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
                         FOREIGN KEY (motorbike_id) REFERENCES users(id) ON DELETE CASCADE,
                         FOREIGN KEY (parcel_id) REFERENCES parcels(id) ON DELETE CASCADE
);
