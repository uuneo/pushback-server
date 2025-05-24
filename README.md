 中文 | **[English](README_EN.md)** | **[日本語](README_JA.md)** | **[한국어](README_KO.md)**
# 反推（Pushback）推送服务适配

> **⚠️ 项目说明**
>
> 本项目参考了 [BARK](https://github.com/Finb/bark-server) —— 一个基于 Golang 的推送服务后端，重写了其接口并采用了 `Gin` 框架，便于扩展维护。
>
> 适用于服务端向 iOS 客户端推送通知，尤其适配 [反推](https://pushback.uuneo.com) App。

---

## ✨ 项目特点

- 基于 Gin 实现，结构清晰
- 支持 Docker 和 Docker Compose 部署
- 提供 SQLite / MySQL 存储支持
- 易于二次开发和集成

---

## 📦 配置说明

`config.yaml` 配置文件必须放置在 `/data/config.yaml`，否则服务将无法启动。

```yaml
system:
  name: "PushbackService"
  user: ""               # 用户名（可选）
  password: ""           # 密码（可选）
  host: "0.0.0.0"        # 服务监听地址
  port: "8080"           # 监听端口，需与 Docker 映射一致
  mode: "release"        # 支持 debug, release, test
  dbType: "default"      # 可选: default (SQLite), mysql
  dbPath: "/data"        # SQLite 文件路径
  maxApnsClientCount: 1  # 最大 APNs 客户端连接数

mysql:                  # 仅在 dbType 为 mysql 时生效
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"

apple:                  # 通常无需修改，仅在自定义构建 App 时需要
  keyId:
  teamId:
  topic:
  apnsPrivateKey:
  adminId:               # 管理员 ID，供 App 获取特殊权限
```

---

## 🛠️ 编译方式

确保你在 Linux 环境中构建，并放置可执行文件与配置文件在同一目录：

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go || echo "编译失败"
```

启动：

```bash
./main
```

---

## 🐳 Docker 部署

你可以使用如下命令快速部署：

```bash
docker run -d --name pushback-server \
  -p 8080:8080 \
  -v ./data:/data \
  --restart=always \
  neouu/pushback:latest
```

或使用 GitHub Container Registry 镜像：

```bash
docker run -d --name pushback-server \
  -p 8080:8080 \
  -v ./data:/data \
  --restart=always \
  ghcr.io/uuneo/pushback:latest
```

---

## 🐳 Docker Compose 部署

1. 将项目中的 `/deploy` 文件夹复制到服务器。
2. 确保 `/data/config.yaml` 存在并配置完整。
3. 启动服务：

```bash
docker-compose up -d
```

---

## 📤 推送示例

以下是一个简单的推送请求示例：

```bash
curl "https://api.day.app/yourkey/测试数据?icon=https://example.com/icon.png&group=提醒"
```

你也可以通过 Telegram 发送以下格式（代码框）：

```
https://api.day.app/yourkey?body=测试数据&image=https://example.com?id=uuid
```

---

## 📎 参考项目

- [BARK](https://github.com/Finb/bark-server) - 推送服务原始实现
- [Gin](https://github.com/gin-gonic/gin) - Web 框架

---

## 📮 联系与支持

如需帮助或想了解更多，可通过 issue 区或邮箱联系维护者。欢迎 PR！
