# 多阶段构建
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o poetry-api main.go

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/poetry-api .
COPY --from=builder /app/config ./config

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./poetry-api"]
