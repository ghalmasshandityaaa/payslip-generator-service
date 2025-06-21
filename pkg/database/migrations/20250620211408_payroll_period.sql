-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "payroll_period" (
    id ulid PRIMARY KEY,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    processed_by ulid DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by ulid NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by ulid
);

ALTER TABLE "payroll_period" ADD CONSTRAINT "fk_payroll_period_created_by" FOREIGN KEY ("created_by") REFERENCES "employee" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "payroll_period" ADD CONSTRAINT "fk_payroll_period_updated_by" FOREIGN KEY ("updated_by") REFERENCES "employee" ("id") ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE "payroll_period" ADD CONSTRAINT "check_payroll_period_dates" CHECK (start_date <= end_date);
ALTER TABLE "payroll_period" ADD CONSTRAINT "check_payroll_period_unique" UNIQUE (start_date, end_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payroll_period;
-- +goose StatementEnd
