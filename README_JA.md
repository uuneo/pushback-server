日本語 | **[中文](README.md)** | **[English](README_EN.md)** | **[한국어](README_KO.md)**
# Pushback 通知サービスアダプター

> **⚠️ プロジェクト説明**
>
> このプロジェクトは [BARK](https://github.com/Finb/bark-server) をベースにしています — Golangベースのプッシュ通知サービスバックエンドです。`Gin`フレームワークを使用してインターフェースを書き直し、保守性と拡張性を向上させました。
>
> サーバーからiOSクライアントへのプッシュ通知に適しており、特に [Pushback](https://pushback.uuneo.com) アプリに最適化されています。

---

## ✨ プロジェクトの特徴

- Ginフレームワークで実装、構造が明確
- DockerとDocker Composeによるデプロイをサポート
- SQLite / MySQLストレージをサポート
- 拡張と統合が容易

---

## 📦 設定ガイド

`config.yaml`設定ファイルは`/data/config.yaml`に配置する必要があります。配置しない場合、サービスは起動しません。

```yaml
system:
  name: "PushbackService"
  user: ""               # ユーザー名（任意）
  password: ""           # パスワード（任意）
  host: "0.0.0.0"        # サービスリッスンアドレス
  port: "8080"           # リッスンポート、Dockerマッピングと一致する必要あり
  mode: "release"        # debug、release、testをサポート
  dbType: "default"      # オプション: default (SQLite)、mysql
  dbPath: "/data"        # SQLiteファイルパス
  maxApnsClientCount: 1  # 最大APNsクライアント接続数

mysql:                  # dbTypeがmysqlの場合のみ有効
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"

apple:                  # 通常は変更不要、カスタムアプリビルド時のみ必要
  keyId:
  teamId:
  topic:
  apnsPrivateKey:
  adminId:               # アプリの特別権限用管理者ID
```

---

## 🛠️ ビルド手順

Linux環境でビルドし、実行ファイルを設定ファイルと同じディレクトリに配置してください：

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go || echo "ビルド失敗"
```

サービスの起動：

```bash
./main
```

---

## 🐳 Dockerデプロイ

以下のコマンドで素早くデプロイできます：

```bash
docker run -d --name pushback-server \
  -p 8080:8080 \
  -v ./data:/data \
  --restart=always \
  neouu/pushback:latest
```

またはGitHub Container Registryイメージを使用：

```bash
docker run -d --name pushback-server \
  -p 8080:8080 \
  -v ./data:/data \
  --restart=always \
  ghcr.io/uuneo/pushback:latest
```

---

## 🐳 Docker Composeデプロイ

1. プロジェクトの`/deploy`フォルダをサーバーにコピーします。
2. `/data/config.yaml`が存在し、適切に設定されていることを確認します。
3. サービスを起動：

```bash
docker-compose up -d
```

---

## 📤 プッシュ通知の例

シンプルなプッシュ通知リクエストの例：

```bash
curl "https://api.day.app/yourkey/テストデータ?icon=https://example.com/icon.png&group=リマインダー"
```

Telegramから以下のフォーマットで送信することもできます（コードブロック）：

```
https://api.day.app/yourkey?body=テストデータ&image=https://example.com?id=uuid
```

---

## 📎 参考プロジェクト

- [BARK](https://github.com/Finb/bark-server) - オリジナルのプッシュ通知サービス実装
- [Gin](https://github.com/gin-gonic/gin) - Webフレームワーク

---

## 📮 お問い合わせとサポート

ヘルプや詳細情報が必要な場合は、issueセクションまたはメールでメンテナーにお問い合わせください。PRも歓迎します！
