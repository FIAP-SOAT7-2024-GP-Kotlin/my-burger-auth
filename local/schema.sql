SET TIME ZONE 'UTC';
CREATE SCHEMA IF NOT EXISTS my_burger;

CREATE TABLE IF NOT EXISTS users
(
    id UUID PRIMARY KEY NOT NULL,
    cpf VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    "role" VARCHAR(255) NOT NULL
);