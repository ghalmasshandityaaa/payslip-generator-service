# Swagger Documentation Setup

This document provides comprehensive information about the Swagger documentation setup for the Payslip Generator Service API.

## Overview

The API documentation is automatically generated using Swagger/OpenAPI 3.0 specifications. The documentation includes:

- **Authentication endpoints** - User sign-in and token management
- **Payroll management** - Period creation, processing, and payslip generation
- **Attendance tracking** - Employee attendance records
- **Overtime management** - Overtime hour tracking and calculations
- **Reimbursement handling** - Expense reimbursement processing

## Features

### 🔐 Authentication
- JWT-based authentication with access and refresh tokens
- Role-based access control (Admin/Employee)
- Secure password validation

### 📊 Payroll System
- Payroll period management
- Automated payroll processing
- Individual payslip generation
- Comprehensive payroll reports

### ⏰ Attendance Tracking
- Work shift recording with start/end times
- Duration calculations
- Daily attendance validation

### 🕒 Overtime Management
- Overtime hour tracking (1-3 hours max per day)
- Overtime pay calculations (2x hourly rate)
- Date-based overtime validation

### 💰 Reimbursement System
- Expense reimbursement requests
- Amount and description tracking
- Reimbursement deduction from salary

## API Endpoints

### Authentication
- `POST /v1/auth/sign-in` - User authentication

### Payroll (Admin & Employee)
- `GET /v1/payroll/period` - List payroll periods
- `POST /v1/payroll/period` - Create payroll period (Admin only)
- `POST /v1/payroll/process` - Process payroll (Admin only)
- `GET /v1/payroll/payslip` - Get employee payslip (Employee only)
- `GET /v1/payroll/payslip/report` - Get payslip report (Admin only)

### Attendance (Employee)
- `POST /v1/attendance` - Create attendance record

### Overtime (Employee)
- `POST /v1/overtime` - Create overtime record

### Reimbursement (Employee)
- `POST /v1/reimbursement` - Create reimbursement record

## Data Models

### Core Entities
- **Employee** - User accounts with roles and salary information
- **PayrollPeriod** - Payroll calculation periods with processing status
- **Attendance** - Work shift records with time tracking
- **Overtime** - Overtime hour records with date validation
- **Reimbursement** - Expense reimbursement requests

### View Models
- **Payslip** - Comprehensive payslip with calculations
- **PayslipReport** - Multi-employee payroll summary
- **WebResponse** - Standardized API response format

## Setup Instructions

### 1. Install Dependencies

```bash
# Install Swagger CLI tools
go install github.com/swaggo/swag/cmd/swag@latest

# Install Fiber Swagger
go get github.com/gofiber/swagger
go get github.com/swaggo/fiber-swagger
```

### 2. Generate Documentation

```bash
# Generate Swagger documentation
make swagger-gen
```

### 3. Start the Server

```bash
# Development mode with hot reload
make run-dev

# Or production mode
make run
```

### 4. Access Documentation

Open your browser and navigate to:
```
http://localhost:8080/swagger/
```

## API Testing

### Authentication Flow

1. **Sign In** - Use the `/v1/auth/sign-in` endpoint with valid credentials
2. **Get Token** - Copy the `access_token` from the response
3. **Authorize** - Click the "Authorize" button in Swagger UI
4. **Enter Token** - Use format: `Bearer <your_access_token>`
5. **Test Endpoints** - All protected endpoints will now work

## Error Handling

The API returns appropriate HTTP status codes:

- `200` - Success
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid/missing token)
- `403` - Forbidden (insufficient permissions)
- `500` - Internal Server Error

## Security

### Authentication
- JWT tokens with configurable expiration
- Refresh token mechanism
- Secure password hashing

### Authorization
- Role-based access control
- Admin vs Employee permissions
- Token validation middleware

### Validation
- Request body validation
- Date/time format validation
- Business rule validation (e.g., overtime hours limit)

## Development

### Adding New Endpoints

1. **Add Swagger annotations** to your handler methods
2. **Update request/response models** with proper documentation
3. **Generate documentation** with `make swagger-gen`
4. **Test the endpoint** in Swagger UI

### Example Handler Annotation

```go
// CreateSomething creates a new resource
// @Summary Create something
// @Description Create a new something with the provided details
// @Tags Something
// @Accept json
// @Produce json
// @Security bearer
// @Param request body model.CreateSomethingRequest true "Something details"
// @Router /something [post]
func (h *Handler) CreateSomething(ctx *fiber.Ctx) error {
    // Implementation here
}
```

### Example Model Annotation

```go
// Something represents a something in the system
// swagger:model Something
type Something struct {
    // Unique identifier
    // example: "01HXYZ123456789ABCDEFGHIJK"
    ID string `json:"id"`
    
    // Name of the something
    // example: "Example Name"
    Name string `json:"name"`
}
```

## Additional Resources

- [Swagger/OpenAPI Specification](https://swagger.io/specification/)
- [Fiber Framework Documentation](https://docs.gofiber.io/)
- [Go Swagger Documentation](https://github.com/swaggo/swag)
