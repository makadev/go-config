# README Documentation Issues

Issues identified by cross-checking README.md against actual source code, tests, and examples.

---

## Issue 1: README struct-as-env-prefix example produces wrong env var names

**Labels:** bug, documentation

**Description:**

The "Struct-as-env-prefix pattern" section (README lines 188–203) shows:

```go
type ServerConfig struct {
    Host string `env:"HOST"` // resolved as SERVER_HOST
    Port int    `env:"PORT"` // resolved as SERVER_PORT
}
type AppConfig struct {
    Server ServerConfig `env:"SERVER"`
}
```

This claims `SERVER_HOST` / `SERVER_PORT`, but the code in `metadata.go` does **not** insert an underscore between the struct env prefix and the nested field's env tag — it concatenates directly:

```go
envVar = env_prefix + envVar   // "SERVER" + "HOST" = "SERVERHOST"
newEnvPrefix = envVar           // newEnvPrefix = "SERVER"
```

All working tests and examples use a trailing underscore in the tag (e.g. `env:"NESTED_"`, `env:"PRIMARY_"`, `env:"SERVER"`→ should be `env:"SERVER_"`).

**Impact:** Users copying the README example will get `SERVERHOST` / `SERVERPORT` instead of the expected `SERVER_HOST` / `SERVER_PORT`.

---

## Issue 2: README shows wrong error message for duplicate config key detection

**Labels:** documentation

**Description:**

The duplicate-detection example (README line 508) shows:

```
invalid config struct: failed to list fields: duplicate config key "host" found
```

Actual error from the code:

- `NewConfig` wraps metadata errors as: `"failed to list fields: <inner error>"`
- The `"invalid config struct:"` prefix is only used for `checkConfigStruct` errors (nil pointer, non-struct), **not** for metadata/duplicate errors.

Real output: `failed to list fields: duplicate config key "host" found`

**Impact:** Misleading for users trying to match or parse error messages programmatically.

---

## Issue 3: README table format examples show space-aligned columns, actual output uses tabs

**Labels:** documentation

**Description:**

The table format examples in the README (lines 300–306, 332–346) show neatly space-aligned columns:

```
CONFIG_KEY   VALUE       SECRET
----------   -----       ------
host         localhost
```

The actual code in `dump.go` uses `\t` (tab characters) as separators:

```go
lines = append(lines, "CONFIG_KEY\tVALUE\tSECRET")
```

Output alignment depends entirely on the terminal's tab-stop settings.

**Impact:** Cosmetic mismatch — actual output will look different from what users expect based on the README.

---

## Issue 4: README documents `metadata` content mode shows ConfigName, but table format omits it

**Labels:** documentation, enhancement

**Description:**

The "Multiple dump content modes" table (README lines 323–328) describes `"metadata"` content as showing:

> Config key + env var name + **config name** per field

While `ConfigName` _is_ populated in the `DumpEntry` struct (and thus appears in JSON/YAML output), the **table format** renderer for `"metadata"` only outputs: `CONFIG_KEY | ENV_VAR | VALUE | SECRET` — **no CONFIG_NAME column**.

**Impact:** The documentation oversells `metadata` mode for table output. Either the table renderer should add a CONFIG_NAME column, or the README should clarify that ConfigName is only available in JSON/YAML formats.

---

## Issue 5: `NewDumpOptions()` defaults to JSON format, but `Dump()` defaults to YAML — undocumented inconsistency

**Labels:** documentation

**Description:**

- `NewDumpOptions()` returns `Format: "json"` as the default.
- `Dump()` hardcodes `Format: "yaml"`.

Neither the README nor the code comments explain this difference. If a user creates options via `NewDumpOptions()` and passes them to `DumpWithOptions()`, they get JSON — not the YAML that `Dump()` produces.

**Impact:** Confusing behavior if a user assumes `NewDumpOptions()` produces the same defaults as `Dump()`.

---

## Issue 6: AutoEnv registers env vars for struct-level fields (not just leaf fields)

**Labels:** documentation, bug

**Description:**

When `AutoEnv = true`, the metadata builder also registers env var entries for struct-typed fields (e.g., a field `Server ServerConfig` with config key `"server"` gets env var `"SERVER"`). The README only mentions leaf field env vars.

If someone sets a `SERVER` env var, `loadFromEnv` would attempt to set a struct-typed field from a plain string value, which would fail with an "unsupported field type" error.

**Impact:** Unexpected errors in production if an env var name happens to collide with a struct-level auto-generated name. The README should document this or the code should skip struct fields during AutoEnv env var generation.

---

## Issue 7: `DumpEnv` README example has no matching struct context

**Labels:** documentation

**Description:**

The `DumpEnv` example (README lines 310–315) shows output like:

```
APP_HOST=localhost
APP_PORT=8080
APP_PASSWORD=***
```

But the env var names (`APP_HOST`, etc.) don't match any struct definition shown nearby in the README. The closest struct uses `env:"APP_HOST"` etc., but it's in a different section.

**Impact:** Minor — readers trying to reproduce the example won't know which struct definition and options produce this exact output. A self-contained example would be clearer.
