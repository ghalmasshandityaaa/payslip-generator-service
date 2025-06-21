# API Documentation

## Overview

The Payslip Generator Service provides a comprehensive REST API for managing employee payroll, attendance, overtime, and reimbursements. This document provides detailed information about all available endpoints, request/response formats, and error handling.

## Base Information

- **Base URL**: `http://localhost:3000/v1`
- **Content-Type**: `application/json`
- **Authentication**: JWT Bearer Token (except for authentication endpoints)
- **API Version**: v1

## Authentication

### JWT Token Format
All authenticated requests must include a JWT token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

### Token Structure
```json
{
  "id": "01JY2PMVA2TGFAB0Y7B2ZPEJST",
  "is_admin": false,
  "exp": 1750495977,
  "iat": 1750492377
}
```

## Common Response Format

### Success Response
```json
{
  "ok": true,
  "data": {
    // Response data
  }
}
```

### Error Response
```json
{
  "ok": false,
  "errors": [
    {
      "field": "username",
      "tag": "required",
      "param": "",
      "value": ""
    }
  ]
}
```

### Pagination Response
```json
{
  "ok": true,
  "data": [
    // Array of items
  ],
  "paging": {
    "page": 1,
    "page_size": 10,
    "total_item": 100,
    "total_page": 10
  }
}
```

## Endpoints

### Authentication

#### POST /auth/sign-in
Authenticate user and receive access tokens.

**Request Body:**
```json
{
  "username": "emp_001",
  "password": "Password123!@#"
}
```

**Response:**
```json
{
  "ok": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid credentials or missing fields
- `401 Unauthorized`: Authentication failed

### Attendance Management

#### POST /attendance
Create a new attendance record for the authenticated employee.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "start_time": "2025-06-18T08:00:00Z",
  "end_time": "2025-06-18T17:00:00Z"
}
```

**Response:**
```json
{
  "ok": true
}
```

**Validation Rules:**
- `start_time`: Required, must be a valid ISO 8601 datetime
- `end_time`: Required, must be a valid ISO 8601 datetime
- `end_time` must be after `start_time`

### Overtime Management

#### POST /overtime
Create a new overtime record for the authenticated employee.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "date": "2025-06-18",
  "total_hours": 2
}
```

**Response:**
```json
{
  "ok": true
}
```

**Validation Rules:**
- `date`: Required, must be a valid date in YYYY-MM-DD format
- `total_hours`: Required, must be a positive integer

### Reimbursement Management

#### POST /reimbursement
Create a new reimbursement request for the authenticated employee.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "amount": 100000,
  "description": "Travel expenses for client meeting"
}
```

**Response:**
```json
{
  "ok": true
}
```

**Validation Rules:**
- `amount`: Required, must be a positive integer
- `description`: Required, must be a non-empty string

### Payroll Management

#### POST /payroll/period
Create a new payroll period (Admin only).

**Headers:**
```
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "start_date": "2025-06-01",
  "end_date": "2025-06-30"
}
```

**Response:**
```json
{
  "ok": true
}
```

**Validation Rules:**
- `start_date`: Required, must be a valid date in YYYY-MM-DD format
- `end_date`: Required, must be a valid date in YYYY-MM-DD format
- `end_date` must be after `start_date`

#### GET /payroll/period
List payroll periods with pagination.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `size` (optional): Items per page (default: 10, max: 100)

**Response:**
```json
{
  "ok": true,
  "data": [
    {
      "id": "01JY7PQRVBZVVMGAN0FQXDQ1B1",
      "start_date": "2025-06-20T00:00:00Z",
      "end_date": "2025-06-21T00:00:00Z",
      "is_generated": false,
      "created_at": "2025-06-21T05:18:21.931273+07:00",
      "created_by": "01JY2PMV9XAB7ZNWDH23D1VJT0",
      "updated_at": "2025-06-21T05:18:21.992274+07:00",
      "updated_by": null
    }
  ],
  "paging": {
    "page": 1,
    "page_size": 10,
    "total_item": 1,
    "total_page": 1
  }
}
```

#### POST /payroll/process
Process payroll for a specific period (Admin only).

**Headers:**
```
Authorization: Bearer <admin_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "period_id": "01JY8V1VHBDSN6YCY707D4P7KR"
}
```

**Response:**
```json
{
  "ok": true
}
```

**Validation Rules:**
- `period_id`: Required, must be a valid ULID

#### GET /payroll/payslip
Get payslip for the authenticated employee for a specific period.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `period_id` (required): Payroll period ID

**Response:**
```json
{
  "ok": true,
  "data": {
    "attendances": [
      {
        "id": "01JY8QQZ1JE7HXDNVTRVXSEFQY",
        "start_time": "2025-06-18T15:00:00+07:00",
        "end_time": "2025-06-18T15:00:01+07:00",
        "created_at": "2025-06-18T15:00:00+07:00",
        "created_by": "01JY2PMVA2TGFAB0Y7B2ZPEJST",
        "updated_at": "2025-06-21T14:55:11.349973+07:00",
        "updated_by": null
      }
    ],
    "overtime": {
      "total_item": 1,
      "total_amount": 100000,
      "total_hours": 3,
      "overtimes": [
        {
          "id": "01JY7H92CPVPVKQPBB1W29Q6RF",
          "date": "2025-06-19T00:00:00Z",
          "total_hours": 3,
          "created_at": "2025-06-21T03:42:57.302554+07:00",
          "created_by": "01JY2PMVA2TGFAB0Y7B2ZPEJST",
          "updated_at": "2025-06-21T03:42:57.302554+07:00",
          "updated_by": null
        }
      ]
    },
    "reimbursement": {
      "total_item": 1,
      "total_amount": 100000,
      "reimbursements": [
        {
          "id": "01JY76M1YPJ41HD3TP0625BCHB",
          "amount": 100000,
          "description": "Reimbursement for travel expenses",
          "created_at": "2025-06-21T00:36:42.966406+07:00",
          "created_by": "01JY2PMV9XAB7ZNWDH23D1VJT0",
          "updated_at": "2025-06-21T00:36:42.966406+07:00",
          "updated_by": null
        }
      ]
    },
    "basic_salary": 3220000,
    "salary": 2146666,
    "take_home_pay": 2149164
  }
}
```

#### GET /payroll/payslip/report
Get comprehensive payroll report for all employees (Admin only).

**Headers:**
```
Authorization: Bearer <admin_token>
```

**Query Parameters:**
- `period_id` (required): Payroll period ID

**Response:**
```json
{
  "ok": true,
  "data": {
    "employees": [
      {
        "id": "01JY2PMVA2TGFAB0Y7B2ZPEJST",
        "username": "emp_001",
        "basic_salary": 3220000,
        "salary": 2146666,
        "take_home_pay": 2149164
      },
      {
        "id": "01JY2PMVA2ZZNC3A2H94K8PX6F",
        "username": "emp_002",
        "basic_salary": 5320000,
        "salary": 0,
        "take_home_pay": 0
      }
    ],
    "total_basic_salary": 557880000,
    "total_salary": 2146666,
    "total_take_home_pay": 2149164
  }
}
```

## Error Handling

### HTTP Status Codes

- `200 OK`: Request successful
- `400 Bad Request`: Invalid request data or validation errors
- `401 Unauthorized`: Authentication required or invalid token
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

### Error Response Format

```json
{
  "ok": false,
  "errors": [
    {
      "field": "username",
      "tag": "required",
      "param": "",
      "value": ""
    },
    {
      "field": "password",
      "tag": "min",
      "param": "8",
      "value": "123"
    }
  ]
}
```

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Limit**: 100 requests per minute per IP address
- **Headers**: Rate limit information is included in response headers:
  - `X-Ratelimit-Limit`: Maximum requests allowed
  - `X-Ratelimit-Remaining`: Remaining requests in current window
  - `X-Ratelimit-Reset`: Time until rate limit resets (seconds)

## Pagination

List endpoints support pagination with the following parameters:

- `page`: Page number (1-based, default: 1)
- `size`: Items per page (default: 10, max: 100)

Pagination metadata is included in the response:

```json
{
  "paging": {
    "page": 1,
    "page_size": 10,
    "total_item": 100,
    "total_page": 10
  }
}
```

## Data Types

### ULID
Universally Unique Lexicographically Sortable Identifier used for primary keys.

### DateTime
ISO 8601 format: `YYYY-MM-DDTHH:mm:ssZ`

### Date
ISO 8601 date format: `YYYY-MM-DD`


## Testing

### Postman Collection
A complete Postman collection is available in `docs/postman/API.postman_collection.json` with:
- Pre-configured requests for all endpoints
- Environment variables for different environments
- Test scripts for automated testing
- Example responses for reference

### Environment Setup
1. Import the Postman collection
2. Set up environment variables:
   - `hostname`: API base URL
   - `version`: API version (v1)
   - `username`: Test username
   - `password`: Test password
3. Run the sign-in request to get a token
4. Use the token for subsequent requests

## Versioning

The API uses URL versioning (`/v1/`) to ensure backward compatibility. When making breaking changes:
1. Create a new version (e.g., `/v2/`)
2. Maintain the old version for a transition period
3. Deprecate the old version after sufficient notice
4. Remove the old version after the deprecation period 