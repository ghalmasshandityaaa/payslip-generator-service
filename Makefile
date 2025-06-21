# Folder where migrations are stored
MIGRATION_DIR := pkg/database/migrations
DSN ?= postgres://user_yc:npg_pCDfEV3vBxK8@ep-frosty-violet-a166yppn.ap-southeast-1.pg.koyeb.app/payslip
GOOSE_OPTIONS := -dir $(MIGRATION_DIR) postgres "$(DSN)"

# Detect OS and set binary extension
UNAME_S := $(shell uname -s 2>/dev/null || echo Windows_NT)
ifeq ($(UNAME_S),Windows_NT)
    BINARY_EXT := .exe
else ifeq ($(findstring MINGW,$(UNAME_S)),MINGW)
    BINARY_EXT := .exe
else ifeq ($(findstring CYGWIN,$(UNAME_S)),CYGWIN)
    BINARY_EXT := .exe
else
    BINARY_EXT :=
endif

# OS-specific commands
ifeq ($(UNAME_S),Windows_NT)
    MKDIR_CMD := if not exist bin mkdir bin
    RM_CMD := if exist bin rmdir /s /q bin
else ifeq ($(findstring MINGW,$(UNAME_S)),MINGW)
    MKDIR_CMD := if not exist bin mkdir bin
    RM_CMD := if exist bin rmdir /s /q bin
else ifeq ($(findstring CYGWIN,$(UNAME_S)),CYGWIN)
    MKDIR_CMD := if not exist bin mkdir bin
    RM_CMD := if exist bin rmdir /s /q bin
else
    MKDIR_CMD := mkdir -p bin
    RM_CMD := rm -rf bin
endif

# Binary path with OS-specific extension
BINARY_PATH := bin/server$(BINARY_EXT)
MAIN_PATH := cmd/main.go

help:
	@echo "Available commands:"
	@echo "[1] migrate-create name=<migration_name>	Create a new migration with the specified name"
	@echo "[2] migrate-up                            	Apply all available migrations"
	@echo "[3] migrate-down                          	Roll back the last migration"
	@echo "[4] migrate-clean                         	Roll back all migrations"
	@echo "[5] migrate-status                        	Show migration status"
	@echo "[6] build                                 	Build the server"
	@echo "[7] rebuild                               	Rebuild the server"
	@echo "[8] run-dev                               	Run the server in development mode(hot reload)"
	@echo "[9] run                                   	Run the server in production mode"
	@echo "[10] clean                                 	Clean built binaries"
	@echo "[11] swagger-gen                           	Generate Swagger documentation"
	@echo "[12] swagger-serve                         	Serve Swagger documentation locally"

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

build:
	@echo "Building server..."
	@$(MKDIR_CMD)
	go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Build completed: $(BINARY_PATH)"

run-dev:
	air

# Check if binary exists, if not build it first
run: $(BINARY_PATH)
	@echo "Starting server..."
	./$(BINARY_PATH)

# This target will build the binary if it doesn't exist or if source files are newer
$(BINARY_PATH): $(MAIN_PATH)
	@echo "Binary not found or outdated, building..."
	@$(MAKE) build

clean:
	@echo "Cleaning built binaries..."
	@$(RM_CMD)
	@echo "Clean completed"

# Force rebuild
rebuild: clean build

# Generate Swagger documentation
swagger-gen:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/main.go -o docs/swagger
	@echo "Swagger documentation generated in docs/"

# Serve Swagger documentation locally
swagger-serve:
	@echo "Serving Swagger documentation at http://localhost:8080/swagger/"
	@echo "Make sure the server is running first with 'make run-dev' or 'make run'"

.PHONY: help migrate-create migrate-up migrate-down migrate-clean migrate-status build run-dev run clean rebuild swagger-gen swagger-serve