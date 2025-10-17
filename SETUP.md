# Crawlly Setup Guide

## Prerequisites

1. Go 1.25+ installed
2. PostgreSQL database running
3. Database credentials configured

## Database Setup

### 1. Create Database

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE crawlly;

# Exit psql
\q
```

### 2. Run Migrations

```bash
# Run the migration SQL file
psql -U postgres -d crawlly -f migrations/001_create_users_table.up.sql
```

### 3. Update Environment Variables

Edit `.env` file with your database credentials:

```
PORT=3000
DATABASE_URL=postgres://postgres:yourpassword@localhost:5432/crawlly?sslmode=disable
ENVIRONMENT=development
LOG_LEVEL=info
```

## Running the Application

```bash
# Install dependencies
go mod download

# Run the application
go run cmd/api/main.go
```

The server will start on port 3000 (or whatever you specified in .env).

## API Endpoints

### Health Check
- **GET** `/health`
- Returns: `OK`

### Register User
- **POST** `/api/auth/register`
- Body:
```json
{
  "email": "user@example.com",
  "name": "John Doe",
  "password": "secure123"
}
```

### Login
- **POST** `/api/auth/login`
- Body:
```json
{
  "email": "user@example.com",
  "password": "secure123"
}
```

## Testing with cURL

See the cURL commands below for testing the API.
