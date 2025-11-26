CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    reservation_id INT REFERENCES reservations(id) ON DELETE CASCADE,
    amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, paid, failed
    payment_method VARCHAR(50),          -- credit_card, gopay, etc.
    midtrans_order_id VARCHAR(100),      -- Order ID dari Midtrans
    va_number VARCHAR(50),
    va_bank VARCHAR(20),             
    payment_time TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);