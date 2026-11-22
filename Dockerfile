# syntax=docker/dockerfile:1

# ---- 构建层 ----
# 跑在构建机原生架构上（--platform=$BUILDPLATFORM），靠 GOOS/GOARCH 交叉编译出目标架构，
# 无需 QEMU 模拟，arm64 也是全速编译
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
# 镜像版本号：CI 里传 tag 名（v1.2.3），本地构建默认 dev
ARG VERSION=dev

WORKDIR /src

# 先只拷依赖清单再下载，代码改动时可命中层缓存
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags "-s -w -X main.Version=${VERSION}" -o /out/BilibiliDanmuRobot ./cli

# ---- 运行层 ----
FROM alpine:3.22

# tzdata 提供时区数据；默认上海时区，可用 -e TZ=Asia/Tokyo 覆盖
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai

COPY --from=builder /out/BilibiliDanmuRobot /app/BilibiliDanmuRobot

# 沿用既有挂载约定：宿主机目录挂到 /app/data，配置放 <宿主>/etc/bilidanmaku-api.yaml
WORKDIR /app/data

ENTRYPOINT ["/app/BilibiliDanmuRobot"]
