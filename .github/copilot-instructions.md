# Copilot Instructions

This is a lightweight Go library (`go-config`). Follow these defaults for all tasks:

- Write idiomatic Go: follow standard naming conventions, keep interfaces small, and prefer composition.
- Match the existing code style — tabs for indentation, no unused imports, errors wrapped with `fmt.Errorf("…: %w", err)`.
- Run `go test -v -race .` to validate changes; use `make coverage` for coverage reports.
- Keep changes minimal and focused. Do not refactor unrelated code.
- For documentation audits use the `readme-audit` custom agent (`.github/agents/readme-audit.agent.md`).
