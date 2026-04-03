# go-config Examples

Runnable examples demonstrating the main features of go-config. Each example is
self-contained with its own `main.go`, config file, and README.

| Example | What it shows |
|---------|---------------|
| [basic](./basic/) | File loading, env overrides, secret masking, dump formats |
| [webserver](./webserver/) | Nested structs, durations, arrays, live HTTP server |
| [database](./database/) | Multiple connections, struct-as-env-prefix, validation |
| [debug](./debug/) | All dump formats and content types, interactive mode |

## Quick start

```bash
cd examples/basic
go run main.go
```

Override values with environment variables:

```bash
APP_NAME="custom" PORT=3000 go run main.go
```

See each example's README for details.