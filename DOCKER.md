# 🐳 Docker Setup for Crawlly

This guide will help you run Crawlly with Docker, which includes everything you need (PostgreSQL + Go app).

## Prerequisites

- Docker installed
- Docker Compose installed

Check if you have them:
```bash
docker --version
docker-compose --version
```

## 🚀 Quick Start

### 1. Start Everything

```bash
docker-compose up -d
```

This will:
- Download PostgreSQL image
- Build your Go application
- Create a network between them
- Run database migrations automatically
- Start both services

### 2. Check Logs

```bash
# View all logs
docker-compose logs -f

# View only app logs
docker-compose logs -f app

# View only database logs
docker-compose logs -f postgres
```

### 3. Test the API

```bash
# Health check
curl http://localhost:3000/health

# Register a user
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "Test User",
    "password": "test123"
  }'

# Login
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "test123"
  }'
```

### 4. Stop Everything

```bash
docker-compose down
```

### 5. Stop and Remove All Data

```bash
# This removes the database volume (deletes all data)
docker-compose down -v
```

---

## 📋 Common Commands

### Rebuild the Application

If you made code changes:
```bash
docker-compose up -d --build
```

### View Running Containers

```bash
docker-compose ps
```

### Access PostgreSQL Database

```bash
docker-compose exec postgres psql -U postgres -d crawlly
```

Inside psql, you can run:
```sql
-- List all tables
\dt

-- See users table structure
\d users

-- Query users
SELECT * FROM users;

-- Exit
\q
```

### Access Application Container Shell

```bash
docker-compose exec app sh
```

### Restart a Service

```bash
# Restart app only
docker-compose restart app

# Restart postgres only
docker-compose restart postgres
```

---

## 🔧 Troubleshooting

### Port Already in Use

If port 3000 or 5432 is already in use, edit `docker-compose.yml`:

```yaml
services:
  app:
    ports:
      - "3001:3000"  # Change 3001 to any available port

  postgres:
    ports:
      - "5433:5432"  # Change 5433 to any available port
```

### Database Connection Issues

Check if PostgreSQL is ready:
```bash
docker-compose logs postgres
```

You should see: "database system is ready to accept connections"

### Application Won't Start

Check application logs:
```bash
docker-compose logs app
```

### Fresh Start

Delete everything and start over:
```bash
docker-compose down -v
docker-compose up -d --build
```

---

## 🎯 What's Included

### Services

1. **PostgreSQL Database**
   - Version: 16 (Alpine)
   - Port: 5432
   - User: postgres
   - Password: postgres
   - Database: crawlly

2. **Go Application**
   - Port: 3000
   - Automatically connects to PostgreSQL
   - Auto-restarts on failure

### Volumes

- `postgres_data` - Persists database data between restarts

### Network

- `crawlly-network` - Private network for app ↔ database communication

---

## 🔐 Default Credentials

**PostgreSQL:**
- Host: localhost (or `postgres` from within Docker network)
- Port: 5432
- User: postgres
- Password: postgres
- Database: crawlly

**Change these in production!**

---

## 📊 Monitoring

### Check Health

```bash
# Application health
curl http://localhost:3000/health

# Database health
docker-compose exec postgres pg_isready -U postgres
```

### Resource Usage

```bash
docker stats
```

---

## 🚀 Production Notes

For production deployment:

1. Change default passwords in `docker-compose.yml`
2. Use environment variables for secrets
3. Enable SSL/TLS for PostgreSQL
4. Use a reverse proxy (Nginx) for HTTPS
5. Set up proper logging and monitoring
6. Use Docker secrets or a secrets manager
7. Configure resource limits

Example production docker-compose:
```yaml
services:
  postgres:
    environment:
      POSTGRES_PASSWORD: ${DB_PASSWORD}  # Use env var
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
```

---

## 📝 Architecture

```
┌─────────────────────────────────────┐
│     Docker Host (Your Machine)      │
│                                     │
│  ┌───────────────────────────────┐ │
│  │   crawlly-network (bridge)    │ │
│  │                               │ │
│  │  ┌──────────┐  ┌───────────┐ │ │
│  │  │   app    │  │ postgres  │ │ │
│  │  │  :3000   │←→│  :5432    │ │ │
│  │  └────┬─────┘  └─────┬─────┘ │ │
│  └───────┼──────────────┼───────┘ │
│          │              │         │
│     Port 3000      Port 5432      │
│          ↓              ↓         │
└──────────┼──────────────┼─────────┘
           │              │
      Your Browser    DB Client
```

---

Happy Crawling! 🕷️
