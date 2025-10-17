# Environment Configuration Guide

## How .env Works with Docker

The `.env` file is **automatically read by Docker Compose** and all variables are available to use in `docker-compose.yml` using the `${VARIABLE_NAME}` syntax.

## Configuration Variables

### Server Settings
```env
PORT=3000
```
- The port your API will run on
- Same for both local and Docker

### Database Settings
```env
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5452
DB_NAME=crawlly
```

**Important Notes:**

1. **DB_HOST:**
   - For **local development**: Use `localhost` or `127.0.0.1`
   - For **Docker**: Automatically changed to `postgres` (the service name) inside containers

2. **DB_PORT:**
   - This is the **external port** (what you connect to from your host machine)
   - Default: `5432` (standard PostgreSQL port)
   - If you changed it to `5452`, that's the port exposed on your host
   - Inside Docker containers, PostgreSQL always runs on internal port `5432`

### How Port Mapping Works

In `docker-compose.yml`:
```yaml
postgres:
  ports:
    - "${DB_PORT}:5432"  # Host:Container
```

This means:
- `DB_PORT` (e.g., 5452) → Exposed on your machine
- `5432` → Internal port inside the container

**Examples:**
- `DB_PORT=5432` → Accessible at `localhost:5432`
- `DB_PORT=5452` → Accessible at `localhost:5452`
- `DB_PORT=9999` → Accessible at `localhost:9999`

## Current Setup

### Your .env file uses:
```env
DB_PORT=5452
```

This means:
- From your **host machine** → Connect to `localhost:5452`
- From **inside Docker** → Connect to `postgres:5432`

### Why the difference?

**From Host (local development):**
```
postgres://postgres:postgres@localhost:5452/crawlly
```
- Uses the external port mapped from Docker

**From Docker (app container):**
```
postgres://postgres:postgres@postgres:5432/crawlly
```
- Uses service name `postgres` as hostname
- Uses internal port `5432`

## Variable Substitution in docker-compose.yml

Docker Compose automatically substitutes variables from `.env`:

```yaml
environment:
  POSTGRES_USER: ${DB_USER}        # → postgres
  POSTGRES_PASSWORD: ${DB_PASSWORD} # → postgres
  POSTGRES_DB: ${DB_NAME}          # → crawlly
ports:
  - "${DB_PORT}:5432"              # → 5452:5432
```

## Changing Database Credentials

### For Docker:

1. Edit `.env`:
```env
DB_USER=myuser
DB_PASSWORD=mysecurepass
DB_NAME=mydb
```

2. Restart containers:
```bash
make clean  # Remove old data
make go     # Start with new credentials
```

### For Production:

**Never commit real passwords to git!**

Use environment variables:
```bash
export DB_PASSWORD=production_password_here
docker-compose up -d
```

Or use a `.env.production` file (add to .gitignore):
```bash
docker-compose --env-file .env.production up -d
```

## Common Scenarios

### Scenario 1: Using Default Port
```env
DB_PORT=5432
```
- Connect from host: `localhost:5432`
- No conflicts if no local PostgreSQL is running

### Scenario 2: Avoiding Port Conflicts
```env
DB_PORT=5452
```
- Connect from host: `localhost:5452`
- Good if you already have PostgreSQL running on 5432

### Scenario 3: Custom Database
```env
DB_USER=crawlly_user
DB_PASSWORD=super_secret_pass
DB_NAME=crawlly_prod
```
- Uses custom credentials
- More secure for production

## Testing Your Configuration

### Check what Docker Compose sees:
```bash
docker-compose config
```
This shows the final configuration with all variables substituted.

### Verify database connection from host:
```bash
psql -h localhost -p 5452 -U postgres -d crawlly
```

### Verify database connection from Docker:
```bash
docker-compose exec app sh
# Inside container:
env | grep DATABASE_URL
```

## Environment-Specific Files

You can create multiple environment files:

- `.env` - Default (git-ignored)
- `.env.example` - Template (committed to git)
- `.env.development` - Development settings
- `.env.production` - Production settings
- `.env.test` - Testing settings

Use specific files:
```bash
docker-compose --env-file .env.production up -d
```

## Security Best Practices

1. **Never commit `.env` to git** (already in .gitignore)
2. **Use strong passwords in production**
3. **Rotate credentials regularly**
4. **Use secrets management in production** (Docker Swarm secrets, Kubernetes secrets, etc.)
5. **Restrict database access** with firewall rules

## Troubleshooting

### "Connection refused" error:
- Check if port is correct
- Verify DB_PORT matches the exposed port
- Check if Docker containers are running: `docker-compose ps`

### "Password authentication failed":
- Verify DB_USER and DB_PASSWORD in .env
- If you changed credentials, you need to clean and restart:
  ```bash
  make clean
  make go
  ```

### "Database does not exist":
- Check DB_NAME in .env
- Verify migration ran: `docker-compose logs postgres`

### Variables not substituting:
- Make sure `.env` is in the same directory as `docker-compose.yml`
- Check syntax: `${VAR_NAME}` not `$VAR_NAME`
- Run `docker-compose config` to see the result
