English | **[中文](README.md)** | **[日本語](README_JA.md)** | **[한국어](README_KO.md)**
# Pushback Notification Service Adapter

> **⚠️ Project Description**
> [DOCS](https://docs.uuneo.com/#/deploy)
> This project is based on [BARK](https://github.com/Finb/bark-server) — a Golang-based push notification service backend. We've rewritten its interfaces using the `Gin` framework for better maintainability and extensibility.
>
> Suitable for server-to-iOS client push notifications, especially optimized for the [Pushback](https://pushback.uuneo.com) App.

---

## ✨ Project Features

- Implemented with Gin framework, clear structure
- Supports Docker and Docker Compose deployment
- Provides SQLite / MySQL storage support
- Easy to extend and integrate
- Built-in health check and monitoring endpoints
- Support for silent push notifications
- Image upload and management capabilities

---

## 📦 Configuration Guide

The configuration file `config.yaml` is optional. The system will look for configuration in the following order:
1. `/data/config.yaml`
2. `./config.yaml`
3. Use system default configuration

```yaml
system: # System Configuration
  name: "pushback" # Service name
  user: "" # Service username
  password: "" # Service password
  address: "0.0.0.0:8080" # Service listening address
  debug: false # Enable debug mode
  dsn: "" # mysql user:password@tcp(host:port)
  maxApnsClientCount: 1 # Maximum APNs client connections

apple: # Apple Push Configuration
  keyId: "BNY5GUGV38" # Key ID
  teamId: "FUWV6U942Q" # Team ID
  topic: "me.uuneo.Meoworld" # Push topic
  develop: false # Development environment
  apnsPrivateKey: |- # APNs private key
    -----BEGIN PRIVATE KEY-----
    MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgvjopbchDpzJNojnc
    o7ErdZQFZM7Qxho6m61gqZuGVRigCgYIKoZIzj0DAQehRANCAAQ8ReU0fBNg+sA+
    ZdDf3w+8FRQxFBKSD/Opt7n3tmtnmnl9Vrtw/nUXX4ldasxA2gErXR4YbEL9Z+uJ
    REJP/5bp
    -----END PRIVATE KEY-----
  adminId: "" # Admin ID
```

---

## 🔌 API Endpoints

### Health & Info
- `GET /ping` or `GET /p` - Basic health check
- `GET /health` or `GET /h` - Health status endpoint
- `GET /info` - System information (version, build date, etc.)

### Device Management
- `POST /register` - Register new device
- `GET /register/:deviceKey` - Register device with key
- `GET /:deviceKey/token` or `GET /:deviceKey/t` - Get device push token

### Push Notifications
- `POST /push` or `POST /p` - Send push notification
- `GET /:deviceKey/:params1/:params2/:params3` - Push with title, subtitle, body
- `GET /:deviceKey/:params1/:params2` - Push with title, body
- `GET /:deviceKey/:params1` - Push with body only
- `GET /:deviceKey/update` or `GET /:deviceKey/u` - Send silent push

### Media Management
- `GET /upload` or `GET /u` - Upload interface
- `POST /upload` or `POST /u` - Upload media
- `GET /image/:filename` or `GET /img/:filename` - Retrieve uploaded image

---

## 🛠️ Installation Guide

### Method 1: Direct Download (Recommended)

Download the latest version for your system from [GitHub Releases](https://github.com/uuneo/pushback-server/releases):

```bash
# Download and extract
wget https://github.com/uuneo/pushback-server/releases/download/2.3.17/pushback-server_2.3.17_linux_amd64.tar.gz
tar -xzf pushback-server_2.3.17_linux_amd64.tar.gz

# Start the service
./pushback-server
```

### Method 2: Build from Source

Ensure you're building in a Linux environment and place the executable file in the same directory as the configuration file:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go || echo "Build failed"
```

Start the service:

```bash
./main
```

---

## 🐳 Docker Deployment

You can quickly deploy using the following command:

```bash
docker run -d --name pushback-server \
  -p 8080:8080 \
  -v ./data:/data \
  --restart=always \
  sunvx/pushback:latest
```

Or use the GitHub Container Registry image:

```bash
docker run -d --name pushback-server \
  -p 8080:8080 \
  -v ./data:/data \
  --restart=always \
  ghcr.io/uuneo/pushback:latest
```

---

## 🐳 Docker Compose Deployment

1. Copy the `/deploy` folder from the project to your server.
2. Ensure `/data/config.yaml` exists and is properly configured.
3. Start the service:

```bash
docker-compose up -d
```

---

## 📎 Reference Projects

- [BARK](https://github.com/Finb/bark-server) - Original push notification service implementation
- [Gin](https://github.com/gin-gonic/gin) - Web framework

---

## 📮 Contact and Support

For help or more information, please contact the maintainers through the issue section or email. PRs are welcome!
