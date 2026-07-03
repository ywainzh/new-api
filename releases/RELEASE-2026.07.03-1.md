# Release 2026.07.03-1

## Summary
- Version/tag: `2026.07.03-1`
- Source branch: `dev`
- Deployment target: `/opt/new-api`
- Image: `ghcr.io/ywainzh/new-api:2026.07.03-1`
- Public host: `https://ywain.zyspeed.xyz`
- Local reverse-proxy target: `http://127.0.0.1:38080`

## Changes
- Fix usage log pages not refreshing to the latest records when clicking search or refreshing the browser page.
- Force usage log list and statistics queries to refetch on mount and bypass frontend GET request deduplication for log reads.
- Keep default log time ranges dynamic unless the user manually edits the time picker, so page refreshes do not stay locked to an old `endTime`.
- Apply the same dynamic default time behavior to common, drawing, and task logs.
- Move usage log query-parameter serialization into a small shared module to remove an internal frontend dependency cycle.

## Config Changes
- `APP_IMAGE` must use the fixed release tag.
- `PORT=3000` remains the container application port.
- `HOST_BIND=127.0.0.1` and `HOST_PORT=38080` expose only the local reverse-proxy port.
- `SESSION_SECRET` must remain a long random value.
- `SQL_DSN`, `LOG_SQL_DSN`, and `REDIS_CONN_STRING` remain unset.

## Data / Migration Notes
- SQLite path: `/opt/new-api/data/new-api.db`
- Migration impact: no schema or data migration.
- Cache impact: container restart clears in-memory frontend/runtime state; no Redis cache is configured for this personal lite deployment.

## Validation
- Local status: `curl http://127.0.0.1:38080/api/status`
- Public HTTPS: `curl -I https://ywain.zyspeed.xyz`
- SSE streaming: verify Nginx still has buffering disabled.
- Deployment dry checks: `git diff --check`
- Frontend checks: `tsgo -b` and `oxlint` on changed usage-log files.
- Runtime usage logs: create or wait for a new log entry, then click search or refresh `/usage-logs/common`; the newest row should appear without switching pages.

## Deployment Steps
1. Confirm this release document is reviewed.
2. Push `dev` and the fixed tag.
3. Wait for GHCR image publication.
4. Deploy from `/opt/new-api` using `DEPLOY_LITE.md`.

## Rollback
- Previous image tag: `2026.07.02-1`
- Rollback command: set `APP_IMAGE` back to `ghcr.io/ywainzh/new-api:2026.07.02-1` and run `docker compose -f docker-compose.lite.yml pull && docker compose -f docker-compose.lite.yml up -d`.

## Risks
- If a user manually selects a fixed time range, the page intentionally keeps using that selected range until reset.
- Browser, reverse proxy, or CDN layers outside this app could still cache static frontend assets briefly, depending on their configuration.
