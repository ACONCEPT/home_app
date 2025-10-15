# Flutter + Golang Login Application

A full-stack login application with a Flutter frontend and Golang backend, running in Docker containers.

## Architecture

- **Frontend**: Flutter web application with a login page
- **Backend**: Golang REST API using Gorilla Mux
- **Containerization**: Docker and Docker Compose for easy deployment
- **Networking**: Bridge network for container communication

## Project Structure

```
.
├── backend/
│   ├── main.go              # Golang API with gorilla/mux
│   ├── go.mod               # Go module dependencies
│   └── Dockerfile           # Backend container image
├── frontend/
│   ├── lib/
│   │   ├── main.dart        # Flutter app entry point
│   │   ├── login_page.dart  # Login UI component
│   │   └── api_service.dart # API client service
│   ├── web/
│   │   ├── index.html       # Web entry point
│   │   └── manifest.json    # Web manifest
│   ├── pubspec.yaml         # Flutter dependencies
│   └── Dockerfile           # Frontend container image
└── docker-compose.yml       # Orchestration configuration
```

## Prerequisites

**For Backend:**
- Docker (version 20.10 or higher)
- Docker Compose (version 2.0 or higher)

**For Frontend (if running locally - recommended):**
- Flutter SDK (version 2.0 or higher)
- Chrome or another web browser

## Getting Started

### Option 1: Run Backend in Docker + Flutter Locally (Recommended)

#### 1. Start the backend

```bash
docker-compose up backend
```

This will start the Golang API server on port 8080.

#### 2. Run the Flutter app locally

In a separate terminal:

```bash
cd frontend
flutter pub get
flutter run -d chrome --web-port=3000
```

The Flutter app will open in Chrome on `http://localhost:3000`.

### Option 2: Run Backend Only (API Testing)

If you just want to test the API:

```bash
docker-compose up backend
```

- **Backend API**: Available at `http://localhost:8080/api`
- Test with curl:
  ```bash
  curl -X POST http://localhost:8080/api/login \
    -H "Content-Type: application/json" \
    -d '{"username":"testuser"}'
  ```

### Option 3: Full Docker Setup (Advanced)

If you want to run both in Docker (requires large Flutter image download):

```bash
docker-compose up --build
```

**Note**: The Flutter Docker image is quite large (~2-3GB) and may take significant time to download on the first build.

### Access Points

- **Frontend** (if running locally): `http://localhost:3000`
- **Frontend** (if running in Docker): `http://localhost:3000`
- **Backend API**: `http://localhost:8080/api`

## Using the Application

1. Open `http://localhost:3000` in your web browser
2. Enter any username in the login form
3. Click the "Login" button
4. You'll receive a success message with a demo token

## API Endpoints

### POST /api/login

Login endpoint that accepts a username.

**Request:**
```json
{
  "username": "john_doe"
}
```

**Response (Success):**
```json
{
  "success": true,
  "message": "Login successful",
  "token": "demo-token-john_doe"
}
```

**Response (Error):**
```json
{
  "success": false,
  "message": "Username is required"
}
```

### GET /api/health

Health check endpoint.

**Response:**
```json
{
  "status": "healthy"
}
```

## Docker Configuration

### Port Mappings

- **Frontend**: Host port `3000` → Container port `80`
- **Backend**: Host port `8080` → Container port `8080`

### Network Configuration

Both containers run on a custom bridge network (`app-network`) which allows:
- Container-to-container communication using service names
- Isolation from other Docker networks
- Port mapping to the host machine

## Development

### Stopping the application

```bash
docker-compose down
```

### Viewing logs

```bash
# All services
docker-compose logs -f

# Backend only
docker-compose logs -f backend

# Frontend only
docker-compose logs -f frontend
```

### Rebuilding after changes

```bash
docker-compose up --build
```

## Troubleshooting

### Frontend can't connect to backend

- Ensure both containers are running: `docker-compose ps`
- Check backend health: `curl http://localhost:8080/api/health`
- Verify CORS headers are set in the backend

### Container build fails

- Clear Docker cache: `docker-compose build --no-cache`
- Check Docker disk space: `docker system df`
- Review build logs for specific errors

## Technology Stack

- **Frontend Framework**: Flutter 3.x
- **Backend Language**: Go 1.21
- **HTTP Router**: Gorilla Mux 1.8.1
- **HTTP Client**: Dart http package
- **Web Server**: Nginx (for serving Flutter web)
- **Container Runtime**: Docker
- **Orchestration**: Docker Compose

## License

This is a demonstration project.
