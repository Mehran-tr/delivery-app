CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    read BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE
);
