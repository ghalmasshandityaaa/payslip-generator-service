-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "attendance" (
    id ulid PRIMARY KEY,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by ulid NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by ulid
);

ALTER TABLE "attendance" ADD CONSTRAINT "fk_attendance_created_by" FOREIGN KEY ("created_by") REFERENCES "employee" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "attendance" ADD CONSTRAINT "fk_attendance_updated_by" FOREIGN KEY ("updated_by") REFERENCES "employee" ("id") ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE "attendance" ADD CONSTRAINT check_attendance_time_order CHECK (end_time > start_time); -- Ensures attendance end time is strictly after start time
ALTER TABLE "attendance" ADD CONSTRAINT check_attendance_same_day CHECK (DATE(end_time) = DATE(start_time)); -- Prevents attendance records from spanning multiple days
ALTER TABLE "attendance" ADD CONSTRAINT check_attendance_no_weekend CHECK (EXTRACT(DOW FROM start_time) NOT IN (0, 6) AND EXTRACT(DOW FROM end_time) NOT IN (0, 6)); -- Restricts attendance records to weekdays only (Monday-Friday)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "attendance";
-- +goose StatementEnd
