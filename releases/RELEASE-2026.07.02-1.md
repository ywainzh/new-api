# Release 2026.07.02-1

## Summary
- Version/tag: `2026.07.02-1`
- Source branch: `dev`
- Deployment target: `/opt/new-api`
- Image: `ghcr.io/ywainzh/new-api:2026.07.02-1`
- Public host: `https://ywain.zyspeed.xyz`
- Local reverse-proxy target: `http://127.0.0.1:38080`

## Changes
- Change the default routing retry count from `0` to `2`, so failed upstream calls can fall back to another channel by default.
- Change the default `codex cli trace` Channel Affinity rule to allow retry after failure.
- Clear Channel Affinity cache after affinity setting changes and channel routing-pool changes.
- When an affinity-selected channel fails with a retryable error, clear the current affinity cache entry and retry normal channel selection from the highest priority layer.
- Update the default frontend settings to show `Retry Times = 2` and the Codex affinity template with `Skip retry on failure` disabled.

## Config Changes
- `APP_IMAGE` must use the fixed release tag.
- `PORT=3000` remains the container application port.
- `HOST_BIND=127.0.0.1` and `HOST_PORT=38080` expose only the local reverse-proxy port.
- `SESSION_SECRET` must remain a long random value.
- `SQL_DSN`, `LOG_SQL_DSN`, and `REDIS_CONN_STRING` remain unset.

## Data / Migration Notes
- SQLite path: `/opt/new-api/data/new-api.db`
- Migration impact: no schema or data migration.
- Cache impact: the personal lite deployment does not use Redis for `new-api`; restarting the container clears in-memory Channel Affinity cache.

## Validation
- Local status: `curl http://127.0.0.1:38080/api/status`
- Public HTTPS: `curl -I https://ywain.zyspeed.xyz`
- SSE streaming: verify Nginx still has buffering disabled.
- Deployment dry checks: `git diff --check`
- Runtime routing: verify failed Codex `/v1/responses` requests no longer stay fixed on a single affinity channel.

## Deployment Steps
1. Confirm this release document is reviewed.
2. Push `dev` and the fixed tag.
3. Wait for GHCR image publication.
4. Deploy from `/opt/new-api` using `DEPLOY_LITE.md`.

## Rollback
- Previous image tag: `2026.06.30-2`
- Rollback command: set `APP_IMAGE` back to `ghcr.io/ywainzh/new-api:2026.06.30-2` and run `docker compose -f docker-compose.lite.yml pull && docker compose -f docker-compose.lite.yml up -d`.

## Risks
- Existing database option overrides for `RetryTimes` or `channel_affinity_setting.*` can supersede these code defaults.
- Requests that explicitly target a specific channel still do not fall back to other channels.
