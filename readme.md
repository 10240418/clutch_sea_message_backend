## 开始
```

docker pull golang:1.24-alpine
docker pull alpine:3.14
docker pull mysql:8.0
$env:DOCKER_BUILDKIT=0; docker build -t hisense-vmi-backend:latest .
docker compose up -d backend
```