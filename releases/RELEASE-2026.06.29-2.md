# Release 2026.06.29-2

## Summary
- Version/tag: `2026.06.29-2`
- Source branch: `dev`
- Deployment target: `/opt/new-api`
- Image: `ghcr.io/ywainzh/new-api:2026.06.29-2`
- Public host: `https://ywain.zyspeed.xyz`
- Local reverse-proxy target: `http://127.0.0.1:38080`

## Changes
- Personal lite deployments now default to the `default` frontend theme at runtime.
- Personal mode keeps wallet/top-up hidden consistently in the mobile drawer.
- Keeps the server-side lite deployment flow unchanged: GHCR builds the image, the server only pulls and restarts it.

## Config Changes
- `APP_IMAGE=ghcr.io/ywainzh/new-api:2026.06.29-2`
- `PORT=3000` remains the container application port.
- `HOST_BIND=127.0.0.1` and `HOST_PORT=38080` expose only `127.0.0.1:38080:3000` on the host.
- `SESSION_SECRET` remains a long random value and must not be replaced with the sample value.
- `PERSONAL_MODE_ENABLED=true`, `SQLITE_PATH=/data/new-api.db`, `MEMORY_CACHE_ENABLED=true`, `ERROR_LOG_ENABLED=false`, `BATCH_UPDATE_ENABLED=false`, and `UPDATE_TASK=false` remain the lite defaults.
- `SQL_DSN`, `LOG_SQL_DSN`, and `REDIS_CONN_STRING` must remain unset.

## Data / Migration Notes
- SQLite path: `/opt/new-api/data/new-api.db`
- Logs path: `/opt/new-api/logs`
- Migration impact: normal application startup migrations only; no external database, Redis, or log database is enabled.
- Existing DB values for `theme.frontend` are normalized to `default` while `PERSONAL_MODE_ENABLED=true`.

## Validation
- Local status: `curl http://127.0.0.1:38080/api/status`
- Public HTTPS: open `https://ywain.zyspeed.xyz`
- Theme check: `/api/status` should report `theme=default`, `version=2026.06.29-2`, and `personal_mode_enabled=true`.
- SSE streaming: confirm streaming responses are not buffered or truncated by Nginx.
- Deployment dry checks: `docker compose -f docker-compose.lite.yml pull`, `docker compose -f docker-compose.lite.yml up -d`, `docker compose -f docker-compose.lite.yml ps`, `docker compose -f docker-compose.lite.yml logs --tail=200 new-api`, and `docker stats --no-stream new-api`.

## Deployment Steps
1. Confirm this release document is reviewed.
2. Push `dev`.
3. Create and push the fixed tag `2026.06.29-2`.
4. Wait for GHCR image publication at `ghcr.io/ywainzh/new-api:2026.06.29-2`.
5. Deploy from `/opt/new-api` using `DEPLOY_LITE.md`.

## Rollback
- Previous image tag: `2026.06.29-1`
- Rollback command: set `APP_IMAGE` back to the previous fixed tag and run `docker compose -f docker-compose.lite.yml pull && docker compose -f docker-compose.lite.yml up -d`.

## Risks
- If the GHCR package is private, the server must be logged in with a token that can read packages.
- If Nginx points to `3000` instead of `38080`, public traffic will fail because the compose file binds only `127.0.0.1:38080`.
- If the browser has cached the old classic assets, a hard refresh may be needed after deployment.
