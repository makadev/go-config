## Known Issues

### Direct field access is not synchronized

Method-level thread safety (`Load` / `Get*` / `Set*` / `Dump*`) is implemented via an internal `sync.RWMutex`, but **direct access to the exported members `Config.Data`, `Config.Metadata`, and `Config.Options` bypasses that mutex**.

Reading or writing these fields concurrently — or mixing direct access with method calls without coordination — is a **data race** and is not safe.

**How to prevent it**

Use the locking helpers that `Config` exposes to share the same mutex:

```go
// Preferred: use the scoped helpers (lock is released via defer)
cfg.WithLock(func() {
    cfg.Data.Counter++          // safe write
})

var snap MyConfig
cfg.WithRLock(func() {
    snap = *cfg.Data            // safe read
})
```

For cases where defer-scoping is not suitable, use the raw methods:

```go
cfg.Lock()
cfg.Data.Counter++
cfg.Unlock()

cfg.RLock()
value := cfg.Data.Counter
cfg.RUnlock()
```

**Deadlock / reentrancy**

`sync.RWMutex` is **not reentrant**. Do **not** call any `Config` method (`Load`, `Dump`, `Get*`, `Set*`, …) while already holding the lock (i.e. inside a `WithLock`/`WithRLock` callback, or between `Lock`/`Unlock` or `RLock`/`RUnlock` calls). Doing so will deadlock immediately.

**`Config.Options` mutations**

`Config.Options` should be set up before the first call to `Load` and left unchanged afterwards. Modifying `Options` fields from multiple goroutines, or while `Load` is running, is not safe even with the helpers above, because `Load` reads `Options` while holding the write lock only for the data update phase — not for the options read phase.

See the *Thread safety* section of [README.md](README.md#thread-safety) for additional examples.
