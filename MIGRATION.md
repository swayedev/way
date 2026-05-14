# Migration Guide for Way

This guide helps you upgrade to Way v1.0.0-rc1 from earlier versions.

## Overview of Changes

Way v1.0.0-rc1 introduces production-readiness improvements:

1. Safe HTTP server defaults (timeouts)
2. Error returns instead of `log.Fatalf` for configuration functions
3. Gated ASCII art logging
4. Comprehensive error path testing

## Migration Steps

### Step 1: Update Error Handling for Configuration Functions

Three functions now return `error` as a second return value:

**Before (pre-rc1):**
```go
// These called log.Fatalf if env vars were missing
encKey := way.getEncryptionKey()
authKey := way.getAuthenticationKey()
storeKey := way.getStoreEncryptionKey()
```

**After (v1.0.0-rc1):**
```go
// These now return errors
encKey, err := way.getEncryptionKey()
if err != nil {
    log.Fatalf("configuration error: %v", err)
    // or handle gracefully in your app
}

authKey, err := way.getAuthenticationKey()
if err != nil {
    log.Fatalf("configuration error: %v", err)
}

storeKey, err := way.getStoreEncryptionKey()
if err != nil {
    log.Fatalf("configuration error: %v", err)
}
```

**Note**: These functions are unexported (lowercase). If you were calling them directly, you'll need to add error handling or use the public `SetSession()` and `InitDBFromConfig()` APIs instead.

### Step 2: Update GenerateRandomKey Error Handling

**Before (pre-rc1):**
```go
// Returned nil on error and logged
key := way.crypto.GenerateRandomKey(32) // could be nil
if key == nil {
    // handle
}
```

**After (v1.0.0-rc1):**
```go
// Now returns error as second value
key, err := way.crypto.GenerateRandomKey(32)
if err != nil {
    log.Fatalf("failed to generate random key: %v", err)
}
```

### Step 3: Review Server Defaults

Way's `New()` function now sets safe HTTP server timeouts by default:

```go
w := way.New()
// w.Server now has:
//   ReadHeaderTimeout: 5s
//   ReadTimeout: 15s
//   WriteTimeout: 15s
//   IdleTimeout: 30s
```

**If you were manually setting timeouts**, verify they are still appropriate:

```go
w := way.New()
// Optional: override defaults if needed
w.Server.ReadTimeout = 30 * time.Second
w.Server.WriteTimeout = 30 * time.Second
```

### Step 4: Enable ASCII Art Logging (Optional)

The startup ASCII art is now disabled by default. To re-enable it:

```bash
export WAY_LOG_ASCII_ART=true
# or set it in your deployment environment
```

### Step 5: Update Dependencies

Ensure your `go.mod` requires Way v1.0.0-rc1 or later:

```bash
go get -u github.com/swayedev/way@v1.0.0-rc1
```

Also update fcrypt to match:

```bash
go get -u github.com/swayedev/fcrypt@v1.0.0-rc1
```

## Backward Compatibility

- All public `Way`, `Context`, `Session`, `DB`, and `crypto` APIs remain stable.
- Breaking changes are limited to:
  - Error returns for configuration functions (if you were calling them directly)
  - `GenerateRandomKey()` signature change
  - HTTP server timeout defaults (review if you had custom timeouts)
  - ASCII art disabled by default

## Testing Your Migration

Run your application's test suite after upgrading:

```bash
go test ./...
go vet ./...
```

If you encounter compilation errors, they are likely related to the three error-returning functions. Add error handling as shown in Step 1 and Step 2.

## Troubleshooting

**Q: My app crashes on startup with "WAY_DEFAULT_COOKIE_ENCRYPTION_KEY is required"**

A: This error is now returned explicitly instead of silently exiting. Either set the environment variable or handle the error in your initialization code:

```go
err := w.InitDBFromConfig() // or your session init code
if err != nil {
    return fmt.Errorf("configuration error: %w", err)
}
```

**Q: Server timeouts are too strict for my use case**

A: Adjust them after calling `way.New()`:

```go
w := way.New()
w.Server.ReadTimeout = 30 * time.Second
w.Server.WriteTimeout = 30 * time.Second
w.Start(":8080")
```

**Q: I want to keep the ASCII art visible**

A: Set the environment variable:

```bash
WAY_LOG_ASCII_ART=true /path/to/app
```

## Questions?

Refer to [API_FREEZE.md](API_FREEZE.md) for the stable API surface and [SECURITY.md](SECURITY.md) for security best practices.
