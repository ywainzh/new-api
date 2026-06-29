# new-api 轻量部署标准文档

本文件用于个人版小服务器部署。服务器只负责拉取固定 tag 镜像、启动容器和健康检查，不在服务器上构建镜像。

## 固定变量

```bash
export DEPLOY_DIR=/opt/new-api
export CONTAINER_NAME=new-api
export IMAGE_REPO=ghcr.io/ywainzh/new-api
export IMAGE_TAG=YYYY.MM.DD-N
export PREV_IMAGE_TAG=
export APP_IMAGE=${IMAGE_REPO}:${IMAGE_TAG}
export PREV_APP_IMAGE=
export PUBLIC_HOST=ywain.zyspeed.xyz
export HOST_BIND=127.0.0.1
export HOST_PORT=38080
```

生产环境禁止使用 `latest`。推送 tag 前，本地仓库必须有对应 `releases/RELEASE-*.md`；服务器部署目录不需要复制 release 文档。

## 运行栈

- 单容器 `new-api`
- SQLite: `/opt/new-api/data/new-api.db`
- 日志目录: `/opt/new-api/logs`
- 容器监听端口: `3000`
- 宿主机绑定: `127.0.0.1:38080`
- 不运行独立 PostgreSQL、MySQL、Redis 或 ClickHouse
- Compose 文件: `docker-compose.lite.yml`
- 服务器部署目录只保留 `data/`、`logs/`、`.env` 和 `docker-compose.lite.yml`

## 服务器准备

```bash
mkdir -p /opt/new-api/data /opt/new-api/logs
cd /opt/new-api
test -f docker-compose.lite.yml
test -f .env
grep '^APP_IMAGE=' .env
grep -q ':latest$' .env && echo 'ERROR: latest is forbidden' && exit 1 || echo 'APP_IMAGE tag ok'
df -h /
docker system df
```

`.env` 必须至少包含：

```bash
APP_IMAGE=ghcr.io/ywainzh/new-api:YYYY.MM.DD-N
PORT=3000
HOST_BIND=127.0.0.1
HOST_PORT=38080
TZ=Asia/Shanghai
SESSION_SECRET=replace_with_a_long_random_string
PERSONAL_MODE_ENABLED=true
SQLITE_PATH=/data/new-api.db
MEMORY_CACHE_ENABLED=true
ERROR_LOG_ENABLED=false
BATCH_UPDATE_ENABLED=false
UPDATE_TASK=false
NODE_NAME=new-api-personal-1
```

不要设置 `SQL_DSN`、`LOG_SQL_DSN`、`REDIS_CONN_STRING`。

`PORT` 是容器内应用监听端口，保持 `3000`。对外只通过宿主机本地端口 `127.0.0.1:38080` 暴露给 Nginx，不直接公开容器端口。

## Nginx 反向代理

将 `ywain.zyspeed.xyz` 的 HTTPS 站点反代到 `http://127.0.0.1:38080`。证书路径按服务器实际签发位置调整。

```nginx
server {
    listen 80;
    server_name ywain.zyspeed.xyz;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    server_name ywain.zyspeed.xyz;

    ssl_certificate /etc/letsencrypt/live/ywain.zyspeed.xyz/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/ywain.zyspeed.xyz/privkey.pem;

    client_max_body_size 128m;

    location / {
        proxy_pass http://127.0.0.1:38080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection $connection_upgrade;
        proxy_buffering off;
        proxy_cache off;
        proxy_read_timeout 3600s;
        proxy_send_timeout 3600s;
    }
}
```

如果 Nginx 未定义 `$connection_upgrade`，在 `http` 块增加：

```nginx
map $http_upgrade $connection_upgrade {
    default upgrade;
    '' close;
}
```

## 首次部署

仅在 release 文档已完成、镜像 tag 已发布、用户明确批准部署后执行：

```bash
cd "$DEPLOY_DIR"
test -d data || mkdir -p data
test -d logs || mkdir -p logs
grep '^APP_IMAGE=' .env
grep -q ':latest$' .env && echo 'ERROR: latest is forbidden' && exit 1 || echo 'APP_IMAGE tag ok'
docker compose -f docker-compose.lite.yml pull && docker compose -f docker-compose.lite.yml up -d
docker compose -f docker-compose.lite.yml ps
docker compose -f docker-compose.lite.yml logs --tail=200 new-api
docker stats --no-stream new-api
docker compose -f docker-compose.lite.yml exec -T new-api wget -qO- http://127.0.0.1:3000/api/status
curl http://127.0.0.1:38080/api/status
```

## 部署后验证

```bash
nginx -t
systemctl reload nginx
curl http://127.0.0.1:38080/api/status
curl -I https://ywain.zyspeed.xyz
```

确认流式请求时，Nginx 不应缓冲或截断 SSE 响应。若前端或 API 流式输出出现一次性吐出、长时间无数据后断开，优先检查站点配置里是否保留了 `proxy_buffering off;`、`proxy_http_version 1.1;` 和较长的 `proxy_read_timeout`。

## 升级部署

```bash
cd "$DEPLOY_DIR"
cp .env ".env.bak.${IMAGE_TAG}"
sed -i "s#^APP_IMAGE=.*#APP_IMAGE=${APP_IMAGE}#" .env
docker compose -f docker-compose.lite.yml pull && docker compose -f docker-compose.lite.yml up -d
docker compose -f docker-compose.lite.yml ps
docker compose -f docker-compose.lite.yml logs --tail=200 new-api
docker stats --no-stream new-api
docker compose -f docker-compose.lite.yml exec -T new-api wget -qO- http://127.0.0.1:3000/api/status
curl http://127.0.0.1:38080/api/status
```

## 旧镜像清理

健康检查成功后，只清理本项目旧镜像，保留当前镜像和上一版回滚镜像。

```bash
df -h /
docker system df

docker image ls --format '{{.Repository}}:{{.Tag}} {{.ID}}' \
  | while read -r image image_id; do
      case "$image" in
        "${IMAGE_REPO}:"*)
          if [ "$image" != "$APP_IMAGE" ] && { [ -z "${PREV_APP_IMAGE:-}" ] || [ "$image" != "$PREV_APP_IMAGE" ]; }; then
            echo "remove old image: $image"
            docker image rm "$image" || true
          fi
          ;;
      esac
    done

df -h /
docker system df
```

禁止用 `docker system prune -a` 作为常规清理命令。

## 回滚

```bash
cd "$DEPLOY_DIR"
sed -i "s#^APP_IMAGE=.*#APP_IMAGE=${IMAGE_REPO}:${PREV_IMAGE_TAG}#" .env
docker compose -f docker-compose.lite.yml pull && docker compose -f docker-compose.lite.yml up -d
docker compose -f docker-compose.lite.yml ps
docker compose -f docker-compose.lite.yml logs --tail=200 new-api
docker compose -f docker-compose.lite.yml exec -T new-api wget -qO- http://127.0.0.1:3000/api/status
curl http://127.0.0.1:38080/api/status
```

## 禁止事项

- 禁止服务器执行 `docker build` 或 `docker compose up -d --build`
- 禁止服务器执行 `bun install`、`bun run build`、`go build`
- 禁止使用 `latest`
- 禁止没有 release 文档就部署
- 禁止健康检查失败后清理上一版回滚镜像
