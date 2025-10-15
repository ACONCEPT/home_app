# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a full-stack authentication application with a **Flutter web frontend** and **Golang backend API**, orchestrated via Docker Compose with PostgreSQL for user storage.

### Key Components

**Backend (`backend/main.go`)**:
- Single-file Go server using Gorilla Mux for routing
- PostgreSQL database with bcrypt password hashing
- CORS enabled for all endpoints (`Access-Control-Allow-Origin: *`)
- All business logic (database, auth, handlers) in one file
- Logging middleware with emojis (✅ ❌ 🚀 📝 💾) for visibility

**Frontend (`frontend/lib/`)**:
- `main.dart`: Entry point with Material 3 theme
- `login_page.dart`: Login UI with username/password + navigation to signup
- `signup_page.dart`: User registration with password confirmation
- `api_service.dart`: HTTP client service (hardcoded to `http://localhost:8080/api`)

**Database**:
- PostgreSQL 15-alpine in Docker with credentials `loginapp:loginapp123`
- Single `users` table: `id`, `username` (unique), `password` (bcrypt hash), `created_at`
- Database connection via environment variable `DATABASE_URL` or falls back to localhost

### API Endpoints

- `POST /api/login` - Accepts `{username, password}`, returns `{success, message, token}`. Also supports legacy username-only auth.
- `POST /api/signup` - Creates user with `{username, password}`, validates duplicates and minimum password length (6 chars)
- `GET /api/health` - Returns `{status: "healthy"}`

## Development Workflow

### Recommended: Run Backend in Docker + Flutter Locally

Start backend with database:
```bash
docker-compose up backend
```

Run Flutter web app in separate terminal:
```bash
cd frontend
flutter pub get
flutter run -d chrome --web-port=3000
```

Access at `http://localhost:3000`

### Backend Only (API Testing)

```bash
docker-compose up backend
```

Test endpoints:
```bash
# Health check
curl http://localhost:8080/api/health

# Signup
curl -X POST http://localhost:8080/api/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

### Full Docker Setup

```bash
docker-compose up --build
```

**Note**: Frontend Docker image is ~2-3GB. Local Flutter development is faster.

## Important Patterns

### Go Dependency Management

After modifying `go.mod` (adding imports), **must** regenerate `go.sum`:
```bash
cd backend
docker run --rm -v "$PWD":/app -w /app golang:1.21-alpine go mod tidy
```

This is critical for Docker builds. The Dockerfile copies both `go.mod` and `go.sum`.

### Database Connection

Backend reads `DATABASE_URL` environment variable. Docker Compose sets this to connect containers:
```
DATABASE_URL=postgresql://loginapp:loginapp123@db:5432/loginapp?sslmode=disable
```

For local Go development outside Docker, defaults to `localhost:5432`.

### Flutter Web Limitations

- `api_service.dart` hardcodes `http://localhost:8080/api` as baseUrl
- Flutter hot reload works with `r` key when running locally
- Browser hard refresh (Cmd+Shift+R) may be needed after significant changes
- `index.html` uses simplified bootstrap: `<script src="flutter_bootstrap.js" async></script>`

### Password Security

- All passwords hashed with bcrypt (cost 14) before database storage
- Never log or return password values
- Login endpoint supports both authenticated (username+password) and legacy (username-only) modes

## Container Architecture

Services communicate via bridge network `app-network`:
- **db**: PostgreSQL on port 5432 (internal), health check with `pg_isready`
- **backend**: Go API on port 8080 (exposed to host), depends on db health
- **frontend**: Nginx serving Flutter web on port 3000 (if running in Docker)

Backend waits for database health before starting (`condition: service_healthy`).

## Viewing Logs

The backend includes comprehensive logging middleware showing:
- HTTP method, URI, remote address, user agent
- Request duration
- Success/failure with emojis

View logs:
```bash
# All services
docker-compose logs -f

# Backend only
docker-compose logs -f backend

# Database only
docker-compose logs -f db
```

## Rebuild and Cleanup

After code changes:
```bash
docker-compose up --build
```

Stop all services:
```bash
docker-compose down
```

Clear all containers and volumes:
```bash
docker-compose down -v
```
