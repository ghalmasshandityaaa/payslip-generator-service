-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "ulid";
CREATE TABLE IF NOT EXISTS "employee" (
    id ulid PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password text NOT NULL,
    salary INTEGER NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE "employee" ADD CONSTRAINT "check_salary_positive" CHECK (salary >= 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "employee";
DROP EXTENSION IF EXISTS "ulid";
-- +goose StatementEnd
