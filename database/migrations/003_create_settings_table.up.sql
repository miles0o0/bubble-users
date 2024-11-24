CREATE TABLE settings (
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    theme VARCHAR(50) DEFAULT 'light',
    notifications BOOLEAN DEFAULT TRUE
);
