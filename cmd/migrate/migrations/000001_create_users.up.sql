CREATE EXTENSION IF NOT EXISTS citex;

CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    email citex UNIQUE NOT NULL,
    username varchar(255) UNIQUE NOT NULL,
    password bytea NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);