CREATE TABLE courts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,           -- Court 1, Court 2, etc.
    status VARCHAR(20) DEFAULT 'active', -- active, maintenance
    location VARCHAR(100) NOT NULL,      -- Location of the court
    price_per_hour DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);