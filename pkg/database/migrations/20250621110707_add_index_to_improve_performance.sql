-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_employee_username ON employee (username);
CREATE INDEX IF NOT EXISTS idx_overtime_date ON overtime (date);
CREATE INDEX IF NOT EXISTS idx_overtime_period ON overtime (date, created_by);
CREATE INDEX IF NOT EXISTS idx_payroll_period_start_date ON payroll_period (start_date);
CREATE INDEX IF NOT EXISTS idx_payroll_period_end_date ON payroll_period (end_date);
CREATE INDEX IF NOT EXISTS idx_payroll_period ON payroll_period (start_date, end_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_employee_username;
DROP INDEX IF EXISTS idx_overtime_date;
DROP INDEX IF EXISTS idx_overtime_period;
DROP INDEX IF EXISTS idx_payroll_period_start_date;
DROP INDEX IF EXISTS idx_payroll_period_end_date;
DROP INDEX IF EXISTS idx_payroll_period;
-- +goose StatementEnd
