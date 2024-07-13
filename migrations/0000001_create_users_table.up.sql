CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    passport_number VARCHAR(20) NOT NULL,
    surname VARCHAR(100),
    name VARCHAR(100),
    patronymic VARCHAR(100),
    address VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

COMMIT;