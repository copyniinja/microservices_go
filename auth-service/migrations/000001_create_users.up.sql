-- +migrate Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,            
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    password TEXT NOT NULL,           
    active BOOLEAN DEFAULT TRUE,     
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
