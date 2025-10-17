# 🔥 Hot Reload Development Guide

## What is Hot Reload?

Hot reload automatically rebuilds and restarts your Go application whenever you save changes to your code. No more manual restarts!

**Without Hot Reload:**
1. Edit code
2. Stop server (Ctrl+C)
3. Run `go run cmd/api/main.go`
4. Wait for compilation
5. Test changes

**With Hot Reload:**
1. Edit code
2. Save file
3. ✅ Done! (Server automatically rebuilds and restarts)

---

## 🚀 Quick Start

### Start Development Mode with Hot Reload

```bash
make dev
```

That's it! Now edit any `.go` file and watch it automatically reload.

### View Logs

```bash
# Logs are shown automatically with make dev
# Or in a separate terminal:
make dev-logs
```

### Stop Development Mode

```bash
# Press Ctrl+C in the terminal running make dev
# Or run:
make dev-down
```

---

## 📋 How It Works

### The Stack

1. **Air** - Go hot reload tool (like nodemon for Node.js)
2. **Docker Volumes** - Syncs your code into the container
3. **PostgreSQL** - Database (same as production mode)

### What Happens When You Save

```
1. You edit internal/membership/handlers/auth_handlers.go
2. Air detects the file change
3. Air rebuilds the application (go build)
4. Air restarts the server
5. Your changes are live! (takes ~1-2 seconds)
```

### Files Being Watched

Air watches all `.go` files in:
- `cmd/`
- `internal/`
- `pkg/`

**Excluded:**
- `*_test.go` (test files)
- `tmp/` (build directory)
- `vendor/`

---

## 🎯 Development Workflow

### Typical Day of Development

```bash
# Morning: Start your environment
make dev

# Work on features, Air auto-reloads on every save
# Edit: internal/membership/usecases/register_user.go
# Save → Auto reload → Test

# Add new endpoints
# Edit: internal/membership/handlers/auth_handlers.go
# Save → Auto reload → Test with Postman

# Database changes? Clean and restart
Ctrl+C
make clean
make dev

# End of day: Stop everything
make dev-down
```

### Testing Your Changes

While `make dev` is running:

```bash
# In another terminal, test your API
curl http://localhost:3000/health

curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "name": "Test", "password": "test123"}'
```

---

## 🔧 Configuration

### .air.toml

Air configuration is in `.air.toml`. Key settings:

```toml
[build]
  cmd = "go build -o ./tmp/main ./cmd/api"  # Build command
  bin = "./tmp/main"                         # Binary location
  include_ext = ["go", "tpl", "tmpl"]       # Watch these extensions
  exclude_dir = ["tmp", "vendor"]           # Ignore these dirs
  delay = 1000                              # Wait 1s before rebuild
```

### Dockerfile.dev

Development Dockerfile includes:
- Full Go environment (not just the binary)
- Air installed
- Source code mounted as volume

### docker-compose.dev.yml

Overrides production settings:
- Uses `Dockerfile.dev`
- Mounts source code as volume
- Runs Air instead of compiled binary

---

## 🆚 Development vs Production Mode

### Development Mode (`make dev`)

```yaml
Dockerfile: Dockerfile.dev
Command: air -c .air.toml
Volumes: Source code mounted
Size: ~300MB (includes Go compiler)
Rebuild: Automatic on file save
Logs: Verbose (debug level)
```

**Use for:**
- Local development
- Writing new features
- Debugging
- Experimenting

### Production Mode (`make go`)

```yaml
Dockerfile: Dockerfile
Command: ./main (binary)
Volumes: None
Size: ~20MB (minimal)
Rebuild: Manual
Logs: Info level
```

**Use for:**
- Testing production build
- Performance testing
- Deployment simulation

---

## 🐛 Troubleshooting

### Hot Reload Not Working

**Problem:** Changes don't trigger rebuild

**Solutions:**
```bash
# 1. Check if Air is running
docker compose -f docker-compose.yml -f docker-compose.dev.yml ps

# 2. Check logs for errors
make dev-logs

# 3. Restart dev environment
make dev-down
make dev
```

### Compilation Errors

**Problem:** Syntax error prevents compilation

**What you'll see:**
```
# command-line-arguments
./main.go:10:2: syntax error: unexpected newline, expecting type
```

**Solution:**
- Fix the syntax error in your code
- Air will automatically retry once you save

### Port Already in Use

**Problem:** Port 3000 is already taken

**Solution:**
```bash
# Stop any other services
make dev-down
make down

# Or change port in .env
PORT=3001
```

### Slow Rebuilds

**Problem:** Takes more than 3-5 seconds to reload

**Possible causes:**
1. Large codebase (normal)
2. Docker resources low
3. Too many dependencies

**Solutions:**
```bash
# Increase Docker resources (Docker Desktop → Settings → Resources)
# Or exclude more directories in .air.toml
```

### Database Connection Lost

**Problem:** App can't connect to database after reload

**Solution:**
```bash
# Database might have stopped
make dev-down
make dev

# Or check database health
docker compose ps
```

---

## 💡 Pro Tips

### Tip 1: Use Multiple Terminals

```bash
# Terminal 1: Run dev mode
make dev

# Terminal 2: Test your API
curl http://localhost:3000/...

# Terminal 3: Database operations
make db-shell
```

### Tip 2: Watch Build Times

Air shows build time in logs:
```
Built successfully in 1.2s
```

If builds are slow, consider:
- Excluding more directories
- Using `go mod vendor`
- Upgrading hardware

### Tip 3: Instant Database Reset

```bash
# Quick way to reset database during development
make clean && make dev
```

### Tip 4: Debug Mode

Edit `docker-compose.dev.yml` to enable more logging:
```yaml
environment:
  LOG_LEVEL: debug
```

### Tip 5: Custom Air Settings

Edit `.air.toml` to customize:
- Build delays
- Excluded directories
- File extensions to watch
- Kill behavior

---

## 📊 Performance Comparison

| Operation | Without Hot Reload | With Hot Reload |
|-----------|-------------------|-----------------|
| Code change | Edit → Stop → Build → Run (10s) | Edit → Save (2s) |
| Add endpoint | Manual restart (10s) | Auto reload (2s) |
| Fix bug | Manual restart (10s) | Auto reload (2s) |
| Test change | Switch terminals + restart | Just save |

**Time saved per day:** ~30-60 minutes (based on 20-40 changes/day)

---

## 🔍 What's Inside

### Files Created for Hot Reload

1. **`.air.toml`**
   - Air configuration
   - Defines what to watch and how to build

2. **`Dockerfile.dev`**
   - Development Docker image
   - Includes Go compiler + Air

3. **`docker-compose.dev.yml`**
   - Development overrides
   - Mounts source code
   - Runs Air

4. **`tmp/`** (auto-created)
   - Air's build directory
   - Git-ignored
   - Contains compiled binaries

5. **`build-errors.log`** (auto-created)
   - Compilation errors
   - Git-ignored
   - Useful for debugging

---

## 🎓 Learning More

### Air Documentation
- GitHub: https://github.com/cosmtrek/air
- Customize: Edit `.air.toml`

### Alternative Tools
- **Fresh**: https://github.com/gravityblast/fresh
- **CompileDaemon**: https://github.com/githubnemo/CompileDaemon
- **Reflex**: https://github.com/cespare/reflex

### Without Docker (Local Development)

Install Air locally:
```bash
go install github.com/cosmtrek/air@latest

# Run
air -c .air.toml
```

---

## 🎉 Summary

**Development Workflow:**
```bash
make dev        # Start with hot reload
# Edit code... save... repeat...
make dev-down   # Stop when done
```

**Production Testing:**
```bash
make go         # Start production mode
make down       # Stop when done
```

**Reset Everything:**
```bash
make clean      # Nuclear option - removes all data
make dev        # Fresh start
```

Happy coding with instant feedback! 🚀
