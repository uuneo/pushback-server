中文 | **[English](README_EN.md)** | **[日本語](README_JA.md)** | **[한국어](README_KO.md)**
# 反推（Pushback）推送服务适配

> **⚠️ 项目说明**
> [DOCS](https://docs.uuneo.com/#/deploy)
> 本项目参考了 [BARK](https://github.com/Finb/bark-server) —— 一个基于 Golang 的推送服务后端，重写了其接口并采用了 `Gin` 框架，便于扩展维护。
>
> 适用于服务端向 iOS 客户端推送通知，尤其适配 [反推](https://pushback.uuneo.com) App。

---

## ✨ 项目特点

- 基于 Gin 实现，结构清晰
- 支持 Docker 和 Docker Compose 部署
- 提供 SQLite / MySQL 存储支持
- 易于二次开发和集成
- 内置健康检查和监控端点
- 支持静默推送通知
- 图片上传和管理功能

---

## 📦 配置说明

配置文件 `config.yaml` 是可选的，系统将按以下顺序查找配置：
1. `/data/config.yaml`
2. `./config.yaml`
3. 使用系统默认配置

```yaml
system: # 系统配置
  name: "pushback" # 服务名称
  user: "" # 服务用户名
  password: "" # 服务密码
  address: "0.0.0.0:8080" # 服务监听地址
  debug: false # 是否开启调试模式
  dsn: "" # mysql user:password@tcp(host:port)
  maxApnsClientCount: 1 # 最大APNs客户端连接数

apple: # 苹果推送配置
  keyId: "BNY5GUGV38" # 密钥ID
  teamId: "FUWV6U942Q" # 团队ID
  topic: "me.uuneo.Meoworld" # 推送主题
  develop: false # 是否开发环境
  apnsPrivateKey: |- # APNs私钥
    -----BEGIN PRIVATE KEY-----
    MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgvjopbchDpzJNojnc
    o7ErdZQFZM7Qxho6m61gqZuGVRigCgYIKoZIzj0DAQehRANCAAQ8ReU0fBNg+sA+
    ZdDf3w+8FRQxFBKSD/Opt7n3tmtnmnl9Vrtw/nUXX4ldasxA2gErXR4YbEL9Z+uJ
    REJP/5bp
    -----END PRIVATE KEY-----
  adminId: "" # 管理员ID
```

### 命令行参数

除了配置文件外，还可以通过命令行参数或环境变量来配置服务：

| 参数 | 环境变量 | 说明 | 默认值 |
|------|----------|------|--------|
| `--addr` | `PB_SERVER_ADDR` | 服务器监听地址 | 空 |
| `--config`, `-c` | `PB_SERVER_CONFIG` | 配置文件路径 | `/data/config.yaml` |
| `--dsn` | `PB_SERVER_DSN` | MySQL DSN | 空 |
| `--maxApnsClientCount`, `-max` | `PB_MAX_APNS_CLIENT_COUNT` | 最大 APNs 客户端数量 | 0（无限制） |
| `--debug` | `PB_DEBUG` | 启用调试模式 | false |
| `--develop`, `-dev` | `PB_DEVELOP` | 启用推送开发模式 | false |
| `--user`, `-u` | `PB_USER` | 服务器用户名 | 空 |
| `--password`, `-p` | `PB_PASSWORD` | 服务器密码 | 空 |

命令行参数优先级高于配置文件，环境变量优先级高于命令行参数。

---

## 🔌 API 接口说明

### 健康检查与信息
- `GET /ping` 或 `GET /p` - 基础健康检查
- `GET /health` 或 `GET /h` - 健康状态端点
- `GET /info` - 系统信息（版本、构建日期等）

### 设备管理
- `POST /register` - 注册新设备
- `GET /register/:deviceKey` - 使用密钥注册设备
- `GET /:deviceKey/token` 或 `GET /:deviceKey/t` - 获取设备推送令牌

### 推送通知
- `POST /push` 或 `POST /p` - 发送推送通知
- `GET /:deviceKey/:params1/:params2/:params3` - 推送标题、副标题、内容
- `GET /:deviceKey/:params1/:params2` - 推送标题、内容
- `GET /:deviceKey/:params1` - 仅推送内容
- `GET /:deviceKey/update` 或 `GET /:deviceKey/u` - 发送静默推送

### 媒体管理
- `GET /upload` 或 `GET /u` - 上传界面
- `POST /upload` 或 `POST /u` - 上传媒体
- `GET /image/:filename` 或 `GET /img/:filename` - 获取已上传图片

---

## 🛠️ 安装方式

### 方式一：直接下载（推荐）

从 [GitHub Releases](https://github.com/uuneo/pushback-server/releases) 下载对应系统的最新版本：

```bash
# 下载并解压
wget https://github.com/uuneo/pushback-server/releases/download/2.3.17/pushback-server_2.3.17_linux_amd64.tar.gz
tar -xzf pushback-server_2.3.17_linux_amd64.tar.gz

# 启动服务
./pushback-server
```

### 方式二：从源码编译

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
  sunvx/pushback:latest
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

* 将项目中的 `/deploy` 文件夹复制到服务器。

```bash
docker-compose up -d
```

---

## 📎 参考项目

- [BARK](https://github.com/Finb/bark-server) - 推送服务原始实现
- [Gin](https://github.com/gin-gonic/gin) - Web 框架

---

## 📮 联系与支持

如需帮助或想了解更多，可通过 issue 区或邮箱联系维护者。欢迎 PR！
