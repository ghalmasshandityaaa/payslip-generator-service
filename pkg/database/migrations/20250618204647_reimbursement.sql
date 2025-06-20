-- +goose Up
-- +goose StatementBegin
CREATE TYPE "reimbursement_status" AS ENUM ('pending', 'approved', 'rejected');
CREATE TABLE IF NOT EXISTS "reimbursement" (
    id ulid PRIMARY KEY,
    amount INTEGER NOT NULL,
    description TEXT NOT NULL,
    status reimbursement_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by ulid NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by ulid
);

ALTER TABLE "reimbursement" ADD CONSTRAINT "fk_reimbursement_created_by" FOREIGN KEY ("created_by") REFERENCES "employee" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "reimbursement" ADD CONSTRAINT "fk_reimbursement_updated_by" FOREIGN KEY ("updated_by") REFERENCES "employee" ("id") ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE "reimbursement" ADD CONSTRAINT "check_reimbursement_amount" CHECK (amount > 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "reimbursement";
DROP TYPE IF EXISTS "reimbursement_status";
-- +goose StatementEnd
