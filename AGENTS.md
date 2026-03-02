# Repository Guidelines

## Project Structure & Module Organization
`cmd/` contains Cobra command implementations grouped by feature (`cmd/issue`, `cmd/workspace`, `cmd/project`). `internal/api/` contains Plane API clients, resolvers, and most test coverage. `internal/config/` manages config and keyring access. `internal/output/` handles table, JSON, and YAML rendering. Shared API types live in `pkg/plane/`. Build entrypoints are `main.go` and the `Makefile`.

## Build, Test, and Development Commands
- `make build`: build the `plane` binary with version metadata.
- `make run`: build and run the CLI locally.
- `make test`: run the full Go test suite.
- `make check`: run `fmt`, `vet`, `lint`, and `test`.
- `make lint`: run `golangci-lint`.
- `make setup-hooks`: install the repo’s `pre-commit` hooks.

For targeted work, prefer package-scoped commands such as `go test ./internal/api -run TestResolveAssigneesIncludesSuggestions`. Live integration tests use build tags and should be serialized to avoid API rate limits:

```bash
go test -tags integration -p 1 -parallel 1 ./internal/api ./cmd/issue ./cmd/workspace -v
```

## Coding Style & Naming Conventions
Use standard Go formatting; run `gofmt` before committing. Keep code ASCII unless the file already requires otherwise. Follow Go naming conventions: exported symbols use `CamelCase`, internal helpers use `camelCase`, and tests follow `TestXxx`. Keep command behavior aligned with the public Plane API; remove stale features rather than preserving undocumented endpoints.

## Testing Guidelines
Use Go’s `testing` package with `testify` assertions. Unit tests live alongside the code as `*_test.go`. Integration tests use the `integration` build tag and rely on `PLANE_API_KEY`, `PLANE_WORKSPACE`, and related env vars. Respect the existing 5-second pacing helper in `internal/integrationtest/ratelimit.go` when adding live API tests.

## Commit & Pull Request Guidelines
Recent history favors short imperative subjects, often with a `fix:` prefix for focused changes, e.g. `fix: Use correct /work-items/ API endpoints`. Keep commits scoped to one concern. PRs should include a clear summary, affected commands/modules, test evidence, and any Plane API doc references used to justify behavior changes.

## Security & Configuration Tips
Never commit API keys or workspace secrets. Use a git-ignored `.env` for local integration work. The CLI stores credentials in the OS keyring and reads user config from `~/.config/plane-cli/config.yaml`.
