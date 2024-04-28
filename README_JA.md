日本語 | **[中文](README.md)** | **[English](README_EN.md)** | **[한국어](README_KO.md)**
# Pushback 通知サービスアダプター

> **⚠️ プロジェクト説明**
> [DOCS](https://docs.uuneo.com/#/deploy)
> このプロジェクトは [BARK](https://github.com/Finb/bark-server) をベースにしています — Golangベースのプッシュ通知サービスバックエンドです。`Gin`フレームワークを使用してインターフェースを書き直し、保守性と拡張性を向上させました。
>
> サーバーからiOSクライアントへのプッシュ通知に適しており、特に [Pushback](https://pushback.uuneo.com) アプリに最適化されています。

---

## ✨ プロジェクトの特徴

- Ginフレームワークで実装、構造が明確
- DockerとDocker Composeによるデプロイをサポート
- SQLite / MySQLストレージをサポート
- 拡張と統合が容易
- 組み込みのヘルスチェックとモニタリングエンドポイント
- サイレントプッシュ通知をサポート
- 画像アップロードと管理機能

---

## 📦 設定ガイド

設定ファイル `config.yaml` は任意です。システムは以下の順序で設定を検索します：
1. `/data/config.yaml`
2. `./config.yaml`
3. システムのデフォルト設定を使用

```yaml
system: # システム設定
  name: "pushback" # サービス名
  user: "" # サービスユーザー名
  password: "" # サービスパスワード
  address: "0.0.0.0:8080" # サービスリッスンアドレス
  debug: false # デバッグモード有効化
  dsn: "" # mysql user:password@tcp(host:port)
  maxApnsClientCount: 1 # 最大APNsクライアント接続数

apple: # アップルプッシュ設定
  keyId: "BNY5GUGV38" # キーID
  teamId: "FUWV6U942Q" # チームID
  topic: "me.uuneo.Meoworld" # プッシュトピック
  develop: false # 開発環境
  apnsPrivateKey: |- # APNs秘密鍵
    -----BEGIN PRIVATE KEY-----
    MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgvjopbchDpzJNojnc
    o7ErdZQFZM7Qxho6m61gqZuGVRigCgYIKoZIzj0DAQehRANCAAQ8ReU0fBNg+sA+
    ZdDf3w+8FRQxFBKSD/Opt7n3tmtnmnl9Vrtw/nUXX4ldasxA2gErXR4YbEL9Z+uJ
    REJP/5bp
    -----END PRIVATE KEY-----
  adminId: "" # 管理者ID
```

---

## 🔌 APIエンドポイント

### ヘルスチェックと情報
- `GET /ping` または `GET /p` - 基本的なヘルスチェック
- `GET /health` または `GET /h` - ヘルスステータスエンドポイント
- `GET /info` - システム情報（バージョン、ビルド日付など）

### デバイス管理
- `POST /register` - 新規デバイス登録
- `GET /register/:deviceKey` - キーを使用したデバイス登録
- `GET /:deviceKey/token` または `GET /:deviceKey/t` - デバイスプッシュトークン取得

### プッシュ通知
- `POST /push` または `POST /p` - プッシュ通知送信
- `GET /:deviceKey/:params1/:params2/:params3` - タイトル、サブタイトル、本文付きプッシュ
- `GET /:deviceKey/:params1/:params2` - タイトル、本文付きプッシュ
- `GET /:deviceKey/:params1` - 本文のみプッシュ
- `GET /:deviceKey/update` または `GET /:deviceKey/u` - サイレントプッシュ送信

### メディア管理
- `GET /upload` または `GET /u` - アップロードインターフェース
- `POST /upload` または `POST /u` - メディアアップロード
- `GET /image/:filename` または `GET /img/:filename` - アップロード済み画像取得

---

## 🛠️ インストール方法

### 方法1：直接ダウンロード（推奨）

[GitHub Releases](https://github.com/uuneo/pushback-server/releases) からお使いのシステムに合わせた最新バージョンをダウンロード：

```bash
# ダウンロードと展開
wget https://github.com/uuneo/pushback-server/releases/download/2.3.17/pushback-server_2.3.17_linux_amd64.tar.gz
tar -xzf pushback-server_2.3.17_linux_amd64.tar.gz

# サービスの起動
./pushback-server
```

### 方法2：ソースからのビルド

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
  sunvx/pushback:latest
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

## 📎 参考プロジェクト

- [BARK](https://github.com/Finb/bark-server) - オリジナルのプッシュ通知サービス実装
- [Gin](https://github.com/gin-gonic/gin) - Webフレームワーク

---

## 📮 お問い合わせとサポート

ヘルプや詳細情報が必要な場合は、issueセクションまたはメールでメンテナーにお問い合わせください。PRも歓迎します！
