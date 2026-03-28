# go-config

A lightweight, generic Go library for loading and managing application configuration. It provides a thin layer of structure on top of standard YAML/JSON file loading, with first-class support for environment variable overrides and user-friendly configuration dumping.

## GOALS
- Provide a robust configuration loading system for Go applications.
- Support loading configuration from various sources: files (YAML/JSON) and environment variables.
- Allow easy mapping of configuration fields to struct fields using tags.
- Facilitate the dumping of configuration data in a user-friendly format.

## NON-GOALS
- Provide a complete solution for all possible configuration scenarios. Instead this should be a robust starting point that can be extended as needed and just adds a little bit of structure and organization on top of YAML/JSON/.. file loading.
- Replace existing configuration management tools or libraries. This library is not intended to be a drop-in replacement for existing tools, but rather a complementary solution specifically with a focus on simplicity and ease of use.
- Verification of configuration values. There are other frameworks that specialize in struct validation f.e. https://github.com/go-playground/validator
- Performance and Memory optimization. A configuration is typically loaded once at startup and then accessed in-memory with a small set of keys/values, so the performance impact is minimal.

## Features

- **Generic `Config[T]` type** — The config container is fully typed. The `Data` field gives direct, type-safe access to your config struct without any casting.
- **File loading (YAML & JSON)** — Loads configuration from the first matching file in a configurable list of paths. Both `.yaml`/`.yml` and `.json` extensions are supported and auto-detected.
- **Priority file lookup** — The default search order is `config.local.yaml → config.local.json → config.yaml → config.json`, so local override files take precedence without any extra code.
- **Environment variable loading** — Fields tagged with `env:"VAR_NAME"` are populated from the matching environment variable after files are loaded, allowing env vars to override file values.
- **Struct tag mapping** — The config key name for a field is resolved from struct tags in priority order: `config` → `yaml` → `json`. If none are present the lowercased field name is used.
- **Secret masking in dumps** — Fields tagged `secret:"true"` have their values replaced with `***` in all dump output by default, preventing accidental secret exposure in logs.
- **Multiple dump formats** — `Dump` / `DumpWithOptions` render configuration as JSON, YAML, plain text (`key=value`), or a tab-separated table.
- **Multiple dump content modes** — Choose what to include in a dump: `"config"` (config keys → values), `"env"` (env var names → values), `"metadata"` (config key + env var + config name), or `"all"` (adds Go struct field path).
- **Getter/Setter API** — `GetFieldValue` / `SetFieldValue` address fields by their Go struct path (e.g. `"Server.Port"`); `GetConfigValue` / `SetConfigValue` address them by their config key (e.g. `"server.port"`).
- **Duplicate detection** — Duplicate config keys or env var names across the struct tree are detected at construction time (`NewConfig`) and returned as an error immediately.

## Undocumented Features

These features are implemented but not yet covered in the README:

- **`AutoEnv` mode** — Setting `Options.AutoEnv = true` automatically derives an environment variable name for every field from its full config key without requiring explicit `env` tags (e.g. `server.port` → `SERVER_PORT`). The `EnvPrefix` option prepends a global prefix to all auto-generated names (e.g. `EnvPrefix = "APP_"` → `APP_SERVER_PORT`).
- **Struct-as-env-prefix pattern** — When `AutoEnv` is `false`, setting `env:"SERVER"` on a struct field causes its tag value to become the prefix for all nested fields' env var names (e.g. a nested field with `env:"HOST"` becomes `SERVER_HOST`).
- **Configurable tag priority** — `Options.ConfigTags` controls which struct tags are checked for the config name and in what order. Defaults to `["config", "yaml", "json"]`.
- **`SkipFiles` / `SkipEnv` flags** — Either loading source can be individually disabled via `Options.SkipFiles` or `Options.SkipEnv`.
- **Nil pointer auto-initialization** — When traversing a struct path through a nil `*struct` pointer (e.g. during `SetFieldValue` or env loading), the pointer is automatically allocated rather than returning an error.
- **Slice env var parsing** — Slice fields (e.g. `[]string`, `[]int`) can be populated from a single comma-separated environment variable value (e.g. `HOSTS=a,b,c`).
- **Map env var parsing** — `map[string]T` fields can be populated from a `KEY1=VAL1,KEY2=VAL2` formatted environment variable value.
- **`time.Duration` env support** — Fields of type `time.Duration` are parsed with `time.ParseDuration`, accepting strings like `"30s"` or `"1h30m"`.
- **Extended boolean parsing** — Boolean fields accept `true/false`, `t/f`, `yes/no`, `y/n`, `1/0`, and `on/off` (case-insensitive) when loaded from environment variables.
- **Direct metadata access** — `cfg.Metadata` is publicly exported and exposes three lookup maps: `FieldPathMap` (by Go field path), `EnvMap` (by env var name), and `KeyMap` (by config key), useful for advanced introspection or tooling.
