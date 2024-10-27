한국어 | **[English](README_EN.md)** | **[日本語](README_JA.md)** | **[中文](README.md)**
# Pushback 알림 서비스 어댑터

> **⚠️ 프로젝트 설명**
>
> 이 프로젝트는 [BARK](https://github.com/Finb/bark-server)를 기반으로 합니다 — Golang 기반의 푸시 알림 서비스 백엔드입니다. `Gin` 프레임워크를 사용하여 인터페이스를 재작성하여 유지보수성과 확장성을 향상시켰습니다.
>
> 서버에서 iOS 클라이언트로의 푸시 알림에 적합하며, 특히 [Pushback](https://pushback.uuneo.com) 앱에 최적화되어 있습니다.

---

## ✨ 프로젝트 특징

- Gin 프레임워크로 구현, 구조가 명확함
- Docker와 Docker Compose 배포 지원
- SQLite / MySQL 스토리지 지원
- 확장 및 통합이 용이함

---

## 📦 설정 가이드

`config.yaml` 설정 파일은 `/data/config.yaml`에 위치해야 합니다. 그렇지 않으면 서비스가 시작되지 않습니다.

```yaml
system:
  name: "PushbackService"
  user: ""               # 사용자 이름 (선택사항)
  password: ""           # 비밀번호 (선택사항)
  host: "0.0.0.0"        # 서비스 리스닝 주소
  port: "8080"           # 리스닝 포트, Docker 매핑과 일치해야 함
  mode: "release"        # debug, release, test 지원
  dbType: "default"      # 옵션: default (SQLite), mysql
  dbPath: "/data"        # SQLite 파일 경로
  maxApnsClientCount: 1  # 최대 APNs 클라이언트 연결 수

mysql:                  # dbType이 mysql인 경우에만 유효
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"

apple:                  # 일반적으로 수정 불필요, 커스텀 앱 빌드 시에만 필요
  keyId:
  teamId:
  topic:
  apnsPrivateKey:
  adminId:               # 앱 특별 권한용 관리자 ID
```

---

## ��️ 빌드 방법

Linux 환경에서 빌드하고 실행 파일을 설정 파일과 같은 디렉토리에 배치하세요:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go || echo "빌드 실패"
```

서비스 시작:

```bash
./main
```

---

## 🐳 Docker 배포

다음 명령어로 빠르게 배포할 수 있습니다:

```bash
docker run -d --name pushback-server \
  -p 8080:8080 \
  -v ./data:/data \
  --restart=always \
  neouu/pushback:latest
```

또는 GitHub Container Registry 이미지 사용:

```bash
docker run -d --name pushback-server \
  -p 8080:8080 \
  -v ./data:/data \
  --restart=always \
  ghcr.io/uuneo/pushback:latest
```

---

## 🐳 Docker Compose 배포

1. 프로젝트의 `/deploy` 폴더를 서버에 복사합니다.
2. `/data/config.yaml`이 존재하고 올바르게 설정되어 있는지 확인합니다.
3. 서비스 시작:

```bash
docker-compose up -d
```

---

## 📤 푸시 알림 예제

간단한 푸시 알림 요청 예제:

```bash
curl "https://api.day.app/yourkey/테스트%20데이터?icon=https://example.com/icon.png&group=알림"
```

Telegram에서 다음 형식으로도 전송할 수 있습니다 (코드 블록):

```
https://api.day.app/yourkey?body=테스트%20데이터&image=https://example.com?id=uuid
```

---

## 📎 참조 프로젝트

- [BARK](https://github.com/Finb/bark-server) - 원본 푸시 알림 서비스 구현
- [Gin](https://github.com/gin-gonic/gin) - 웹 프레임워크

---

## 📮 문의 및 지원

도움이나 추가 정보가 필요하시면 issue 섹션이나 이메일로 관리자에게 문의하세요. PR도 환영합니다!
