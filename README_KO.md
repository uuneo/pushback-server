한국어 | **[English](README_EN.md)** | **[日本語](README_JA.md)** | **[中文](README.md)**
# Pushback 알림 서비스 어댑터

> **⚠️ 프로젝트 설명**
> [DOCS](https://docs.uuneo.com/#/deploy)
> 이 프로젝트는 [BARK](https://github.com/Finb/bark-server)를 기반으로 합니다 — Golang 기반의 푸시 알림 서비스 백엔드입니다. `Gin` 프레임워크를 사용하여 인터페이스를 재작성하여 유지보수성과 확장성을 향상시켰습니다.
>
> 서버에서 iOS 클라이언트로의 푸시 알림에 적합하며, 특히 [Pushback](https://pushback.uuneo.com) 앱에 최적화되어 있습니다.

---

## ✨ 프로젝트 특징

- Gin 프레임워크로 구현, 구조가 명확함
- Docker와 Docker Compose 배포 지원
- SQLite / MySQL 스토리지 지원
- 확장 및 통합이 용이함
- 내장된 헬스 체크 및 모니터링 엔드포인트
- 무음 푸시 알림 지원
- 이미지 업로드 및 관리 기능

---

## 📦 설정 가이드

설정 파일 `config.yaml`은 선택사항입니다. 시스템은 다음 순서로 설정을 검색합니다:
1. `/data/config.yaml`
2. `./config.yaml`
3. 시스템 기본 설정 사용

```yaml
system: # 시스템 설정
  name: "pushback" # 서비스 이름
  user: "" # 서비스 사용자 이름
  password: "" # 서비스 비밀번호
  address: "0.0.0.0:8080" # 서비스 리스닝 주소
  debug: false # 디버그 모드 활성화
  dsn: "" # mysql user:password@tcp(host:port)
  maxApnsClientCount: 1 # 최대 APNs 클라이언트 연결 수

apple: # 애플 푸시 설정
  keyId: "BNY5GUGV38" # 키 ID
  teamId: "FUWV6U942Q" # 팀 ID
  topic: "me.uuneo.Meoworld" # 푸시 토픽
  develop: false # 개발 환경
  apnsPrivateKey: |- # APNs 개인 키
    -----BEGIN PRIVATE KEY-----
    MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgvjopbchDpzJNojnc
    o7ErdZQFZM7Qxho6m61gqZuGVRigCgYIKoZIzj0DAQehRANCAAQ8ReU0fBNg+sA+
    ZdDf3w+8FRQxFBKSD/Opt7n3tmtnmnl9Vrtw/nUXX4ldasxA2gErXR4YbEL9Z+uJ
    REJP/5bp
    -----END PRIVATE KEY-----
  adminId: "" # 관리자 ID
```

---

## 🔌 API 엔드포인트

### 헬스 체크 및 정보
- `GET /ping` 또는 `GET /p` - 기본 헬스 체크
- `GET /health` 또는 `GET /h` - 헬스 상태 엔드포인트
- `GET /info` - 시스템 정보 (버전, 빌드 날짜 등)

### 디바이스 관리
- `POST /register` - 새 디바이스 등록
- `GET /register/:deviceKey` - 키를 사용한 디바이스 등록
- `GET /:deviceKey/token` 또는 `GET /:deviceKey/t` - 디바이스 푸시 토큰 가져오기

### 푸시 알림
- `POST /push` 또는 `POST /p` - 푸시 알림 전송
- `GET /:deviceKey/:params1/:params2/:params3` - 제목, 부제목, 본문이 있는 푸시
- `GET /:deviceKey/:params1/:params2` - 제목, 본문이 있는 푸시
- `GET /:deviceKey/:params1` - 본문만 있는 푸시
- `GET /:deviceKey/update` 또는 `GET /:deviceKey/u` - 무음 푸시 전송

### 미디어 관리
- `GET /upload` 또는 `GET /u` - 업로드 인터페이스
- `POST /upload` 또는 `POST /u` - 미디어 업로드
- `GET /image/:filename` 또는 `GET /img/:filename` - 업로드된 이미지 가져오기

---

## 🛠️ 설치 방법

### 방법 1: 직접 다운로드 (권장)

[GitHub Releases](https://github.com/uuneo/pushback-server/releases)에서 시스템에 맞는 최신 버전을 다운로드하세요:

```bash
# 다운로드 및 압축 해제
wget https://github.com/uuneo/pushback-server/releases/download/2.3.17/pushback-server_2.3.17_linux_amd64.tar.gz
tar -xzf pushback-server_2.3.17_linux_amd64.tar.gz

# 서비스 시작
./pushback-server
```

### 방법 2: 소스에서 빌드

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
  sunvx/pushback:latest
```

또는 GitHub Container Registry 이미지를 사용:

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

## 📎 참조 프로젝트

- [BARK](https://github.com/Finb/bark-server) - 원본 푸시 알림 서비스 구현
- [Gin](https://github.com/gin-gonic/gin) - 웹 프레임워크

---

## 📮 문의 및 지원

도움이나 추가 정보가 필요하시면 issue 섹션이나 이메일로 관리자에게 문의하세요. PR도 환영합니다!
