# Setup and Deployment Guide

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Local Development Setup](#local-development-setup)

## Prerequisites

### Required Software

- **Go 1.24.2 or higher**
  ```bash
  # Download from https://golang.org/dl/
  # Or use package manager
  # Ubuntu/Debian
  sudo apt-get install golang-go
  
  # macOS
  brew install go
  
  # Windows
  # Download installer from https://golang.org/dl/
  ```

- **PostgreSQL 12 or higher**
  ```bash
  # Ubuntu/Debian
  sudo apt-get install postgresql postgresql-contrib
  
  # macOS
  brew install postgresql
  
  # Windows
  # Download from https://www.postgresql.org/download/windows/
  ```

- **Make** (for build automation)
  ```bash
  # Ubuntu/Debian
  sudo apt-get install make
  
  # macOS
  # Usually pre-installed
  
  # Windows
  # Install via Chocolatey: choco install make
  # Or use WSL
  ```

### Optional Software

- **Docker** (for containerized deployment)
- **Air** (for hot reload during development)
- **Goose** (for database migrations)

## Local Development Setup

### Step 1: Clone the Repository

```bash
git clone <repository-url>
cd payslip-generator-service
```

### Step 2: Install Dependencies

```bash
go mod download
```

### Step 3: Install Development Tools

```bash
# Install Air for hot reload
go install github.com/cosmtrek/air@latest

# Install Goose for database migrations
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Step 4: Configure Environment

```bash
# Copy configuration template
cp config/config.json.example config/config.json

# Edit configuration for local development
nano config/config.json
```

**Local Configuration Example:**
```json
{
    "App": {
        "Name": "Payslip Generator Service",
        "Version": "1.0.0",
        "Domain": "localhost",
        "Host": "localhost",
        "Port": 3000,
        "Env": "Development",
        "Debug": true,
        "ReadTimeout": 5,
        "WriteTimeout": 5,
        "Prefork": false,
        "SSL": false
    },
    "Security": {
        "CORS": {
            "AllowedOrigins": "http://localhost:3000,http://localhost:8080",
            "AllowedMethods": "GET,POST,PUT,DELETE,OPTIONS",
            "AllowCredentials": true
        },
        "CSRF": {
            "Enabled": false
        },
        "JWT": {
            "Issuer": "http://localhost:3000",
            "Audience": "http://localhost:3000",
            "Subject": "payslip_generator_service",
            "SigningMethod": "HS256",
            "AccessTokenLifetime": 720,
            "AccessTokenSecret": "dev-secret-key-change-in-production",
            "RefreshTokenLifetime": 720,
            "RefreshTokenSecret": "dev-refresh-secret-change-in-production"
        },
        "RateLimit": {
            "Duration": 60,
            "MaxRequests": 1000
        },
        "Crypto": {
            "Key": "dev-crypto-key-change-in-production"
        }
    },
    "Logger": {
        "Level": 6,
        "Pretty": true
    },
    "Postgres": {
        "ConnMaxIdleTime": 30,
        "ConnMaxLifetime": 3600,
        "MaxIdleCons": 5,
        "MaxOpenCons": 10,
        "User": "postgres",
        "Password": "password",
        "Host": "localhost",
        "Port": 5432,
        "Dbname": "payslip_dev",
        "Driver": "postgres",
        "SSLMode": "disable",
        "DryRun": false
    }
}
```

### Step 5: Set Up Database

```bash
# Create database
createdb payslip_dev

# Run migrations
make migrate-up

# Verify migration status
make migrate-status
```

### Step 6: Run the Application

```bash
# Development mode with hot reload
make run-dev

# Or production mode
make run
```

### Step 7: Verify Installation

```bash
# Test the API
curl http://localhost:3000/v1/auth/sign-in \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"emp_001","password":"Password123!@#"}'
```
