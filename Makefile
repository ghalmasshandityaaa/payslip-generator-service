# Folder where migrations are stored
MIGRATION_DIR := pkg/database/migrations
DSN ?= postgres://user_yc:npg_pCDfEV3vBxK8@ep-frosty-violet-a166yppn.ap-southeast-1.pg.koyeb.app/payslip
GOOSE_OPTIONS := -dir $(MIGRATION_DIR) postgres "$(DSN)"

help:
	@echo "Available commands:"
	@echo "[1] migrate-create name=<migration_name>	Create a new migration with the specified name"
	@echo "[2] migrate-up                            	Apply all available migrations"
	@echo "[3] migrate-down                          	Roll back the last migration"
	@echo "[4] migrate-clean                         	Roll back all migrations"
	@echo "[5] migrate-status                        	Show migration status"

migrate-create:
ifdef name
	goose -dir $(MIGRATION_DIR) create $(name) sql
else
	@echo "Error: Please provide a migration name. Usage: make migrate-create name=<migration_name>"
endif

migrate-up:
	goose $(GOOSE_OPTIONS) up

migrate-down:
	goose $(GOOSE_OPTIONS) down

migrate-clean:
	goose $(GOOSE_OPTIONS) down-to 0

migrate-status:
	goose $(GOOSE_OPTIONS) status

