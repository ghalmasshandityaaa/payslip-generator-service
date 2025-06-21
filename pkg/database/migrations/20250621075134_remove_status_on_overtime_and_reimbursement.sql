-- +goose Up
-- +goose StatementBegin
ALTER TABLE "overtime" DROP COLUMN "status";
ALTER TABLE "reimbursement" DROP COLUMN "status";

DROP TYPE IF EXISTS "overtime_status";
DROP TYPE IF EXISTS "reimbursement_status";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TYPE "overtime_status" AS ENUM ('pending', 'approved', 'rejected');
CREATE TYPE "reimbursement_status" AS ENUM ('pending', 'approved', 'rejected');

ALTER TABLE "overtime" ADD COLUMN "status" overtime_status NOT NULL DEFAULT 'pending';
ALTER TABLE "reimbursement" ADD COLUMN "status" reimbursement_status NOT NULL DEFAULT 'pending';
-- +goose StatementEnd
