# Branching Guide for Personal Development

This fork is maintained with two long-lived branches:

- `main`: upstream sync branch for `QuantumNous/new-api`. Keep this branch clean and do not place personal changes here.
- `dev`: personal development and deployment branch for `ywainzh`. Put all local customization, personal features, and deployment changes here.

## Remote Layout

- `origin`: `git@github-ywainzh:ywainzh/new-api.git`
- `upstream`: `https://github.com/QuantumNous/new-api.git`

Use `origin` for this fork. Use `upstream` only to fetch changes from the original project.

## Daily Development

Work on `dev`:

```powershell
git switch dev
git status --short --branch
```

Before making changes, confirm the working tree is clean or understand any existing local changes. Do not move personal work onto `main`.

## Syncing Upstream

When the upstream project updates, first update `main` from `upstream/main`:

```powershell
git switch main
git fetch upstream
git merge --ff-only upstream/main
git push origin main
```

Then bring those updates into `dev`:

```powershell
git switch dev
git merge main
git push origin dev
```

If conflicts appear while merging into `dev`, preserve the user's personal changes unless the user explicitly asks to discard them. Keep project branding and protected attributions intact according to `AGENTS.md`.

## Deployment

Deploy from `dev`, not `main`. Treat `main` as a clean mirror of upstream history and `dev` as the user's runnable personal version.
