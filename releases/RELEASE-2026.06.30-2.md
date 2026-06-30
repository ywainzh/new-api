# Release 2026.06.30-2

## Summary
- Version/tag: `2026.06.30-2`
- Source branch: `dev`
- Deployment target: `/opt/new-api`
- Image: `ghcr.io/ywainzh/new-api:2026.06.30-2`
- Public host: `https://ywain.zyspeed.xyz`
- Local reverse-proxy target: `http://127.0.0.1:38080`

## Changes
- Align default OpenAI pricing for the deployed channel models with official pricing.
- Add default token pricing for `gpt-5.5`, `gpt-5.4`, `gpt-5.4-mini`, `gpt-5.4-nano`, and `gpt-5.3-codex`.
- Add `gpt-image-2` pricing, including text input, cached input, image input, and output ratios.
- Treat `gpt-image-2-4k` as the same pricing alias as `gpt-image-2`.
- Leave non-OpenAI models and provider-specific pricing untouched.

## Config Changes
- `APP_IMAGE` must use the fixed release tag.
- `PORT=3000` remains the container application port.
- `HOST_BIND=127.0.0.1` and `HOST_PORT=38080` expose only the local reverse-proxy port.
- `SESSION_SECRET` must remain a long random value.
- `SQL_DSN`, `LOG_SQL_DSN`, and `REDIS_CONN_STRING` remain unset.

## Data / Migration Notes
- SQLite path: `/opt/new-api/data/new-api.db`
- Migration impact: no schema or data migration.
- Pricing impact: existing usage logs keep their recorded historical pricing metadata; new requests use the corrected defaults after deployment.

## Validation
- Local status: `curl http://127.0.0.1:38080/api/status`
- Public HTTPS: `curl -I https://ywain.zyspeed.xyz`
- SSE streaming: verify Nginx still has buffering disabled.
- Deployment dry checks: `git diff --check`

## Deployment Steps
1. Confirm this release document is reviewed.
2. Push `dev` and the fixed tag.
3. Wait for GHCR image publication.
4. Deploy from `/opt/new-api` using `DEPLOY_LITE.md`.

## Rollback
- Previous image tag: `2026.06.30-1`
- Rollback command: set `APP_IMAGE` back to `ghcr.io/ywainzh/new-api:2026.06.30-1` and run `docker compose -f docker-compose.lite.yml pull && docker compose -f docker-compose.lite.yml up -d`.

## Risks
- `gpt-image-2-4k` is a channel/provider alias rather than a separate official OpenAI model name, so it intentionally follows `gpt-image-2` pricing.
- Direct database option overrides for model ratios, if added later, can supersede these code defaults.
