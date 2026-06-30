# Release 2026.06.30-1

## Summary
- Version/tag: `2026.06.30-1`
- Source branch: `dev`
- Deployment target: `/opt/new-api`
- Image: `ghcr.io/ywainzh/new-api:2026.06.30-1`
- Public host: `https://ywain.zyspeed.xyz`
- Local reverse-proxy target: `http://127.0.0.1:38080`

## Changes
- Channel management now shows used quota only and removes remaining-balance query UI.
- Models can be synced from configured channel abilities into the model catalog.
- New channels trigger one non-blocking model metadata sync after creation.
- Added a model page action to manually sync channel models later.
- Unknown upstream models are still created as basic enabled model records.

## Config Changes
- `APP_IMAGE` must use the fixed release tag.
- `PORT=3000` remains the container application port.
- `HOST_BIND=127.0.0.1` and `HOST_PORT=38080` expose only the local reverse-proxy port.
- `SESSION_SECRET` must be replaced with a long random value.
- `SQL_DSN`, `LOG_SQL_DSN`, and `REDIS_CONN_STRING` remain unset.

## Data / Migration Notes
- SQLite path: `/opt/new-api/data/new-api.db`
- Migration impact: no schema migration is required.
- Runtime behavior: model/vendor metadata rows may be created from channel abilities during channel creation or manual sync.

## Validation
- Local status: `curl http://127.0.0.1:38080/api/status`
- Public HTTPS: `curl -I https://ywain.zyspeed.xyz`
- SSE streaming: confirm streaming responses are not buffered by Nginx.
- Deployment dry checks: `git diff --check`, i18n sync report clean.

## Deployment Steps
1. Confirm this release document is reviewed.
2. Push `dev` and the fixed tag.
3. Wait for GHCR image publication.
4. Deploy from `/opt/new-api` using `DEPLOY_LITE.md`.

## Rollback
- Previous image tag: `2026.06.29-2`
- Rollback command: set `APP_IMAGE` back to `ghcr.io/ywainzh/new-api:2026.06.29-2` and run `docker compose -f docker-compose.lite.yml pull && docker compose -f docker-compose.lite.yml up -d`.

## Risks
- The first sync may create model metadata records for all enabled channel abilities that are missing from the model catalog.
- Official metadata source failures do not block sync; missing models are created as basic records instead.
- Local Go/Bun checks were not available in this workstation environment, so the release relies on repository static checks and the GHCR Docker build.
