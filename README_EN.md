English | **[中文](README.md)** | **[日本語](README_JA.md)** | **[한국어](README_KO.md)**
# Pushback Notification Service Adapter

> **⚠️ Project Description**
>
> This project is based on [BARK](https://github.com/Finb/bark-server) — a Golang-based push notification service backend. We've rewritten its interfaces using the `Gin` framework for better maintainability and extensibility.
>
> Suitable for server-to-iOS client push notifications, especially optimized for the [Pushback](https://your-app-link) App.

---

## ✨ Project Features

- Implemented with Gin framework, clear structure
- Supports Docker and Docker Compose deployment
- Provides SQLite / MySQL storage support
- Easy to extend and integrate

---

## 📦 Configuration Guide

The `config.yaml` configuration file must be placed at `/data/config.yaml`, otherwise the service will not start.

```yaml
system:
  name: "PushbackService"
  user: ""               # Username (optional)
  password: ""           # Password (optional)
  host: "0.0.0.0"        # Service listening address
  port: "8080"           # Listening port, must match Docker mapping
  mode: "release"        # Supports debug, release, test
  dbType: "default"      # Options: default (SQLite), mysql
  dbPath: "/data"        # SQLite file path
  maxApnsClientCount: 1  # Maximum APNs client connections

mysql:                  # Only effective when dbType is mysql
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"

apple:                  # Usually no need to modify, only required for custom App builds
  keyId:
  teamId:
  topic:
  apnsPrivateKey:
  adminId:               # Admin ID for App special permissions
```

---

## 🛠️ Build Instructions

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
  neouu/pushback:latest
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

## 📤 Push Notification Examples

Here's a simple push notification request example:

```bash
curl "https://api.day.app/yourkey/Test%20Data?icon=https://example.com/icon.png&group=Reminder"
```

You can also send via Telegram using the following format (code block):

```
https://api.day.app/yourkey?body=Test%20Data&image=https://example.com?id=uuid
```

---

## 📎 Reference Projects

- [BARK](https://github.com/Finb/bark-server) - Original push notification service implementation
- [Gin](https://github.com/gin-gonic/gin) - Web framework

---

## 📮 Contact and Support

For help or more information, please contact the maintainers through the issue section or email. PRs are welcome!
