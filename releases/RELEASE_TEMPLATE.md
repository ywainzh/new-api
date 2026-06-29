# Release YYYY.MM.DD-N

## Summary
- Version/tag:
- Source branch: `dev`
- Deployment target: `/opt/new-api`
- Image: `ghcr.io/ywainzh/new-api:YYYY.MM.DD-N`
- Public host:
- Local reverse-proxy target: `http://127.0.0.1:38080`

## Changes
- 

## Config Changes
- `APP_IMAGE` must use the fixed release tag.
- `PORT=3000` remains the container application port.
- `HOST_BIND=127.0.0.1` and `HOST_PORT=38080` expose only the local reverse-proxy port.
- `SESSION_SECRET` must be replaced with a long random value.
- `SQL_DSN`, `LOG_SQL_DSN`, and `REDIS_CONN_STRING` remain unset.

## Data / Migration Notes
- SQLite path: `/opt/new-api/data/new-api.db`
- Migration impact:

## Validation
- Local status: `curl http://127.0.0.1:38080/api/status`
- Public HTTPS:
- SSE streaming:
- Deployment dry checks:

## Deployment Steps
1. Confirm this release document is reviewed.
2. Push `dev` and the fixed tag.
3. Wait for GHCR image publication.
4. Deploy from `/opt/new-api` using `DEPLOY_LITE.md`.

## Rollback
- Previous image tag:
- Rollback command: set `APP_IMAGE` back to the previous fixed tag and run `docker compose -f docker-compose.lite.yml pull && docker compose -f docker-compose.lite.yml up -d`.

## Risks
- 
