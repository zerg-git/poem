# 部署指南

## 部署方式概览

本项目支持多种部署方式：

1. **Docker Compose** - 推荐用于快速部署
2. **传统部署** - 适用于已有服务器环境
3. **云服务部署** - 适用于云平台

---

## 一、Docker Compose部署（推荐）

### 1.1 准备工作

确保已安装：
- Docker 20.10+
- Docker Compose 2.0+

### 1.2 后端Dockerfile

创建 `backend/Dockerfile`:

```dockerfile
# 多阶段构建
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o poetry-api main.go

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/poetry-api .
COPY --from=builder /app/config ./config

EXPOSE 8080

CMD ["./poetry-api"]
```

### 1.3 前端Dockerfile

创建 `frontend/Dockerfile`:

```dockerfile
# 构建阶段
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# 运行阶段 - 使用nginx
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### 1.4 Nginx配置

创建 `frontend/nginx.conf`:

```nginx
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;
    index index.html;

    # 前端路由
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Gzip压缩
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
}
```

### 1.5 Docker Compose配置

创建 `docker-compose.yml`:

```yaml
version: '3.8'

services:
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: poetry-backend
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - ./chinese-poetry-master:/app/data:ro
      - ./poems.db:/app/poems.db
    environment:
      - GIN_MODE=release
      - DATA_PATH=/app/data
      - DB_PATH=/app/poems.db
      - PORT=8080
    networks:
      - poetry-network

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: poetry-frontend
    restart: unless-stopped
    ports:
      - "80:80"
    depends_on:
      - backend
    networks:
      - poetry-network

networks:
  poetry-network:
    driver: bridge
```

### 1.6 启动服务

```bash
# 构建并启动
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 1.7 访问应用

- 前端: http://localhost
- 后端API: http://localhost:8080/api/v1

---

## 二、传统部署方式

### 2.1 服务器要求

- Linux服务器 (Ubuntu 20.04+ 推荐)
- 2GB+ 内存
- 10GB+ 磁盘空间
- Go 1.21+
- Node.js 16+
- Nginx

### 2.2 后端部署

```bash
# 1. 上传后端代码到服务器
scp -r backend/ user@server:/var/www/poem/backend

# 2. SSH登录服务器
ssh user@server

# 3. 进入后端目录
cd /var/www/poem/backend

# 4. 安装依赖
go mod download

# 5. 编译
go build -o poetry-api main.go

# 6. 创建systemd服务
sudo nano /etc/systemd/system/poetry-api.service
```

**systemd服务配置**:

```ini
[Unit]
Description=Poetry API Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/poem/backend
ExecStart=/var/www/poem/backend/poetry-api
Restart=always
RestartSec=5
Environment="GIN_MODE=release"
Environment="DATA_PATH=/var/www/poem/chinese-poetry-master"
Environment="PORT=8080"

[Install]
WantedBy=multi-user.target
```

```bash
# 7. 启动服务
sudo systemctl daemon-reload
sudo systemctl enable poetry-api
sudo systemctl start poetry-api
sudo systemctl status poetry-api
```

### 2.3 前端部署

```bash
# 1. 上传前端代码
scp -r frontend/ user@server:/var/www/poem/frontend

# 2. SSH登录服务器
ssh user@server

# 3. 进入前端目录
cd /var/www/poem/frontend

# 4. 安装依赖并构建
npm install
npm run build

# 5. 配置Nginx
sudo nano /etc/nginx/sites-available/poetry
```

**Nginx配置**:

```nginx
server {
    listen 80;
    server_name poetry.example.com;
    root /var/www/poem/frontend/dist;
    index index.html;

    # 前端路由
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API代理
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Gzip压缩
    gzip on;
    gzip_types text/plain text/css application/json application/javascript;
}
```

```bash
# 6. 启用站点
sudo ln -s /etc/nginx/sites-available/poetry /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## 三、云服务部署

### 3.1 阿里云/腾讯云部署

**ECS云服务器部署**:
- 参考传统部署方式
- 注意配置安全组开放80、8080端口

**使用容器服务**:
- 推荐使用阿里云容器服务 ACK
- 或腾讯云容器服务 TKE

### 3.2 Vercel + Railway部署

**前端部署到Vercel**:

```bash
# 安装Vercel CLI
npm i -g vercel

# 部署
cd frontend
vercel
```

**后端部署到Railway**:

1. 连接GitHub仓库
2. 选择backend目录
3. 配置环境变量:
   - `DATA_PATH`: `/workspace/chinese-poetry`
   - `PORT`: `8080`

### 3.3 AWS部署

**使用AWS ECS**:

```yaml
# aws-task-definition.json
{
  "family": "poetry-app",
  "containerDefinitions": [
    {
      "name": "frontend",
      "image": "...",
      "portMappings": [{"containerPort": 80}],
      "memory": 256
    },
    {
      "name": "backend",
      "image": "...",
      "portMappings": [{"containerPort": 8080}],
      "memory": 512,
      "environment": [
        {"name": "DATA_PATH", "value": "/app/data"}
      ]
    }
  ]
}
```

---

## 四、HTTPS配置

### 4.1 使用Let's Encrypt

```bash
# 安装certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d poetry.example.com

# 自动续期
sudo certbot renew --dry-run
```

### 4.2 Nginx HTTPS配置

```nginx
server {
    listen 443 ssl http2;
    server_name poetry.example.com;

    ssl_certificate /etc/letsencrypt/live/poetry.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/poetry.example.com/privkey.pem;

    # 其他配置...
}

server {
    listen 80;
    server_name poetry.example.com;
    return 301 https://$server_name$request_uri;
}
```

---

## 五、监控与日志

### 5.1 日志管理

```bash
# 查看后端日志
sudo journalctl -u poetry-api -f

# 查看Nginx日志
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### 5.2 性能监控

使用Prometheus + Grafana：

```yaml
# docker-compose.monitoring.yml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana
    ports:
      - "3001:3000"
```

---

## 六、备份与恢复

### 6.1 数据备份

```bash
# 备份脚本
#!/bin/bash
DATE=$(date +%Y%m%d)
tar -czf poetry-backup-$DATE.tar.gz chinese-poetry-master/
scp poetry-backup-$DATE.tar.gz backup-server:/backups/
```

### 6.2 自动备份

```bash
# 添加到crontab
crontab -e

# 每天凌晨2点备份
0 2 * * * /path/to/backup-script.sh
```

---

## 七、故障排查

### 常见问题

| 问题 | 可能原因 | 解决方案 |
|------|----------|----------|
| 502 Bad Gateway | 后端服务未启动 | 检查systemd服务状态 |
| CORS错误 | Nginx配置错误 | 检查proxy_set_header配置 |
| 静态资源404 | 路径配置错误 | 检查root路径是否正确 |
| 数据加载失败 | 数据路径错误 | 检查DATA_PATH环境变量 |

### 健康检查脚本

```bash
#!/bin/bash
# health-check.sh

# 检查后端
curl -f http://localhost:8080/api/v1/dynasties || echo "Backend down"

# 检查前端
curl -f http://localhost/ || echo "Frontend down"
```
