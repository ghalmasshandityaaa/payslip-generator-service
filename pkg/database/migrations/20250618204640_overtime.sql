-- +goose Up
-- +goose StatementBegin
CREATE TYPE "overtime_status" AS ENUM ('pending', 'approved', 'rejected');
CREATE TABLE IF NOT EXISTS "overtime" (
    id ulid PRIMARY KEY,
    date DATE NOT NULL,
    total_hours INTEGER NOT NULL,
    status overtime_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by ulid NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by ulid
);

ALTER TABLE "overtime" ADD CONSTRAINT "fk_overtime_created_by" FOREIGN KEY ("created_by") REFERENCES "employee" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "overtime" ADD CONSTRAINT "fk_overtime_updated_by" FOREIGN KEY ("updated_by") REFERENCES "employee" ("id") ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE "overtime" ADD CONSTRAINT "check_overtime_duration"  CHECK (total_hours >= 1 AND total_hours <= 3); -- min 1 hour and max 3 hours is the maximum overtime duration

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "overtime";
DROP TYPE IF EXISTS "overtime_status";
-- +goose StatementEnd
