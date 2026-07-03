# syntax=docker/dockerfile:1

# ---------- 构建阶段 ----------
FROM golang:1.24-alpine AS builder
WORKDIR /src
# 先拷依赖清单并下载,利用层缓存(源码变动不必重下依赖)
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# 纯静态编译:MySQL 驱动为纯 Go,关掉 CGO 可跨镜像即插即用
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath -ldflags "-s -w" \
    -o /out/pickcampus-backend cmd/app/main.go

# ---------- 运行阶段 ----------
FROM alpine:3.20
# ca-certificates: 走 TLS 的外部依赖需要; tzdata: DSN loc=Local 需要时区库
RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -u 10001 pickcampus
WORKDIR /app
COPY --from=builder /out/pickcampus-backend /app/pickcampus-backend
# 基线配置(非敏感默认值);环境相关值与密钥在 compose 里用环境变量覆盖
COPY configs/config.yaml.example /app/configs/config.yaml
USER pickcampus
EXPOSE 8080
ENTRYPOINT ["/app/pickcampus-backend", "-c", "/app/configs/config.yaml"]
