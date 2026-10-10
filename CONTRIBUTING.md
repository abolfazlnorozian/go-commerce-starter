# Contributing

Thanks for interest in improving this starter kit. Contributions that keep the project honest, modular, and easy to run locally are welcome.

## Before you start

1. Open an issue for larger changes so direction can be discussed early.
2. Keep scope focused — prefer small PRs over large rewrites.
3. Do not commit secrets (`.env`, API keys, merchant IDs, production dumps).

## Development setup

```bash
git clone https://github.com/abolfazlnorozian/go-commerce-starter.git
cd go-commerce-starter
cp .env.example .env
# edit .env
go run main.go
```

Optional: `docker-compose up -d --build`

## Branching

- `main` — stable default branch
- Feature / fix branches: `feat/<short-name>`, `fix/<short-name>`, `docs/<short-name>`

## Commit conventions

Prefer concise, imperative subjects (Conventional Commits style is encouraged):

- `feat: add coupon expiry check`
- `fix: prevent nil cart on empty user`
- `docs: clarify docker mongo URL`
- `chore: bump CI go version`

Explain *why* in the body when the change is non-obvious.

## Pull requests

1. Fork (or branch) and implement against latest `main`.
2. Ensure the project builds: `go build ./...`
3. Run tests (and add tests for new changes where applicable): ‍‍‍‍‍‍‍‍`go test ./...`
4. Update docs when behavior or env vars change.
5. Open a PR with a short summary and test notes.

CI (`.github/workflows/go.yml`) runs build and test on pushes and PRs to `main`.

## Code guidelines

- Match existing package layout (`router` / `adminroute` → services → `entities` / `database`).
- Prefer env vars over hardcoded credentials or environment-specific URLs.
- Avoid drive-by refactors unrelated to the PR.
- New behavior should include tests where practical (coverage is currently thin — adding tests is especially valuable).

## Reporting security issues

Do not file public issues for vulnerabilities. Follow [SECURITY.md](./SECURITY.md).
