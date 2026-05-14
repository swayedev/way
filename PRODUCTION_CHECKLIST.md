# Production Checklist for Way v1.0.0-rc1

Use this checklist before deploying a Way application or tagging a Way release candidate.

## Framework Configuration

- Use `way.New()` so the default HTTP server timeouts are applied.
- If you replace the server with `SetServer()`, configure `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, and `IdleTimeout`.
- Keep recovery, request IDs, CORS, security headers, authentication, authorization, and rate limiting as explicit middleware for your application.
- Leave `WAY_LOG_ASCII_ART` unset in production unless you intentionally want startup art in logs.

## Crypto And Sessions

- Way crypto helpers delegate to fcrypt `v1.0.0-rc1` and preserve hex string ciphertext output.
- Use fcrypt directly when you need secret-manager backed key rotation or key identifiers.
- Set session and secure-cookie keys from a secret manager or secure environment variables.
- Prefer `StoreE`, `CookieE`, `DefaultSessionE`, and `DefaultCookieE` for error-returning session lookups.
- Use `HttpOnly`, `Secure`, and appropriate `SameSite` settings for application cookies.

## Database

- Import only the database driver adapter packages your application needs.
- Use parameterized queries and never concatenate user input into SQL.
- Keep database credentials out of source code and logs.
- Verify missing or unsupported drivers fail fast with clear startup errors.

## Logging

- Confirm logs do not include SQL args, full query values, headers, cookies, tokens, passwords, or DSNs.
- Treat request paths and route parameters as potentially sensitive when designing application middleware.

## Outbound HTTP

- `ProxyMedia` uses a timeout-bound HTTP client by default.
- Set a custom client with `SetHTTPClient()` when proxying through controlled networks or when shorter timeouts are required.

## Release Checks

```bash
GOCACHE=/tmp/fileserver-go-cache go test ./...
go vet ./...
```

Before tagging, verify:

- `VERSION` is `1.0.0-rc1`.
- README and API docs mention fcrypt `v1.0.0-rc1`.
- Public API names in `API_FREEZE.md` match the implementation.
- No default logs expose secrets or request credential material.
