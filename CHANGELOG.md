# Changelog

All notable changes to Way will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0-rc1] – 2026-05-13

### Added

- **Safe HTTP Server Defaults**: `way.New()` now sets production-safe timeouts on the HTTP server:
  - `ReadHeaderTimeout: 5s` to prevent slow-read DoS attacks
  - `ReadTimeout: 15s` for request reading
  - `WriteTimeout: 15s` for response writing
  - `IdleTimeout: 30s` for keep-alive idle connections
- **Error Returns for Configuration**: Functions `getEncryptionKey()`, `getAuthenticationKey()`, and `getStoreEncryptionKey()` now return `error` instead of calling `log.Fatalf`, allowing graceful error handling.
- **GenerateRandomKey Error Handling**: `way/crypto.GenerateRandomKey()` now returns `([]byte, error)` instead of `([]byte)`, enabling proper error propagation.
- **Comprehensive Test Coverage**: Added unit tests for error paths, server defaults, and crypto functions.
- **API Freeze Documentation**: Added [API_FREEZE.md](API_FREEZE.md) to document the stable public API surface.
- **Security Policy**: Added [SECURITY.md](SECURITY.md) with security considerations, best practices, and hardening guidance.
- **Migration Guide**: Added [MIGRATION.md](MIGRATION.md) to help users upgrade from pre-1.0 versions.
- **Linter Configuration**: Added `.golangci.yml` for code quality checks.

### Changed

- **ASCII Art Logging**: Server startup ASCII art is now disabled by default and can be enabled with `WAY_LOG_ASCII_ART=true` environment variable.
- **README Updates**: Updated [README.md](README.md) to reflect v1.0.0-rc1 status and link documentation.
- **VERSION File**: Updated to `1.0.0-rc1`.

### Breaking Changes

1. **Error Returns**: `getEncryptionKey()`, `getAuthenticationKey()`, and `getStoreEncryptionKey()` now return `([]byte, error)` instead of `[]byte`. Update your code to handle errors.
2. **GenerateRandomKey Signature**: `way/crypto.GenerateRandomKey()` now returns `([]byte, error)` instead of `[]byte`.
3. **HTTP Server Timeouts**: `way.New()` now sets default timeouts. If you relied on different timeout behavior, update your server configuration after calling `New()`.
4. **ASCII Art Disabled**: Startup ASCII art is now off by default; set `WAY_LOG_ASCII_ART=true` to re-enable.

### Security

- Addressed security concerns with timeout configuration and error handling.
- Improved session and cookie management with explicit error returns.
- Ensured compliance with Go security best practices.

### Documentation

- Added comprehensive API freeze, security, and migration documentation.
- Updated README to reflect production-ready status.

---

## [0.3.0-rc0] – Earlier Development Version

(Previous development versions not tracked in detail. Upgrade to 1.0.0-rc1 for production use.)
