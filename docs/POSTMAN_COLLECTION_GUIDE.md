# Postman Collection Guide

## Overview

This guide provides detailed instructions for using the Postman collection to test and interact with the Payslip Generator Service API. The collection includes all available endpoints with pre-configured requests, environment variables, and test scripts.

## Collection Structure

### Available Collections

1. **API.postman_collection.json** - Main API collection
2. **Local - Admin.postman_environment.json** - Admin user environment
3. **Local - Employee.postman_environment.json** - Employee user environment

### Collection Organization

```
API Collection
├── Auth
│   └── sign-in
├── Attendance
│   └── Create Attendance
├── Overtime
│   └── Create Overtime
├── Reimbursement
│   └── Create Reimbursement
└── Payroll
    ├── Create Payroll Period
    ├── List Payroll Periods
    ├── Process Payroll Period
    ├── Get Payslip
    └── Get Payslip Report
```

## Setup Instructions

### Online Access

The Postman collection can also be accessed online through the following URL:

https://www.postman.com/planetary-meadow-507576/payslip-generator-service/collection/d9hi9vp/api

This allows you to:
- View and fork the collection without downloading
- Run requests directly from your browser
- Collaborate with team members easily
- Stay updated with the latest API changes

To use the online collection:
1. Click the URL above
2. Click "Fork Collection" to create your own copy
3. Configure your environment variables
4. Start making API requests


### Step 1: Import Collections

1. **Open Postman**
2. **Import Collections**:
   - Click "Import" button
   - Select `docs/postman/API.postman_collection.json`
   - Click "Import"

3. **Import Environments**:
   - Import `docs/postman/Local - Admin.postman_environment.json`
   - Import `docs/postman/Local - Employee.postman_environment.json`

### Step 2: Configure Environment Variables

#### Admin Environment Variables
```json
{
  "hostname": "http://localhost:3000",
  "version": "v1",
  "username": "admin_user",
  "password": "Admin123!@#",
  "token": ""
}
```

#### Employee Environment Variables
```json
{
  "hostname": "http://localhost:3000",
  "version": "v1",
  "username": "emp_001",
  "password": "Password123!@#",
  "token": ""
}
```

### Step 3: Select Environment

1. **Choose Environment**: Select either "Local - Admin" or "Local - Employee" from the environment dropdown
2. **Verify Variables**: Check that all variables are properly set

## Authentication Flow

### Step 1: Sign In

1. **Navigate to**: `Auth > sign-in`
2. **Method**: POST
3. **URL**: `{{hostname}}/{{version}}/auth/sign-in`
4. **Headers**:
   ```
   Content-Type: application/json
   ```
5. **Body** (raw JSON):
   ```json
   {
     "username": "{{username}}",
     "password": "{{password}}"
   }
   ```

### Step 2: Extract Token

The sign-in request includes a test script that automatically extracts the access token:

```javascript
pm.test("Status code is 200", function() {
    pm.response.to.have.status(200);
    if (pm.response.status === "OK") {
        const body = pm.response.json();
        pm.environment.set("token", body.data.access_token);
    }
})
```

### Step 3: Use Token

All subsequent requests automatically use the token via the Authorization header:
```
Authorization: Bearer {{token}}
```
