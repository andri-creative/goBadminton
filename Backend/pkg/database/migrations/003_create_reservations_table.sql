CREATE TABLE reservations (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    court_id INT REFERENCES courts(id) ON DELETE CASCADE,
    reservation_date DATE NOT NULL,      -- Tanggal booking
    time_slot VARCHAR(20) NOT NULL,     -- 10:00-11:00, 11:00-12:00, etc.
    duration_hours INT DEFAULT 1,        -- NEW: Durasi dalam jam
    total_amount DECIMAL(10,2) NOT NULL, -- NEW: Total harga
    status VARCHAR(20) DEFAULT 'pending', -- pending, confirmed, cancelled
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);