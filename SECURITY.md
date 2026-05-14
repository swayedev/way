# Security Policy for Way

## Reporting a Vulnerability

If you discover a security vulnerability in Way, please **do not** open a public GitHub issue. Instead:

1. Email the maintainer(s) with a detailed description of the vulnerability.
2. Include steps to reproduce (if applicable).
3. Allow up to 90 days for a fix and coordinated disclosure before public announcement.

## Security Considerations

Way is a lightweight web framework designed for building Go applications. The following security practices are recommended:

### Cryptography

- **Passphrase-Based Encryption**: Way delegates cryptographic operations to the `fcrypt` package (`github.com/swayedev/fcrypt` v1.0.0-rc1).
  - Encryption is performed using AES-GCM, which provides authenticated encryption.
  - Key derivation uses scrypt with cryptographically secure parameters.
  - Way's compatibility crypto helpers return hex strings and accept passphrases. Applications that need secret-manager integration, key identifiers, or rotation should use fcrypt's production APIs directly.

- **Session and Cookie Security**:
  - Sessions use `gorilla/sessions` with either `CookieStore` or custom `SecureCookie` implementations.
  - Encryption keys and authentication keys should be:
    - Generated securely (e.g., using `crypto/rand` or `way/crypto.GenerateRandomKey()`).
    - Of sufficient length (minimum 32 bytes recommended).
    - Managed securely and rotated regularly.
    - Never committed to version control.
  - Set `WAY_DEFAULT_STORE_ENCRYPTION_KEY`, `WAY_DEFAULT_COOKIE_ENCRYPTION_KEY`, and `WAY_DEFAULT_COOKIE_AUTHENTICATION_KEY` environment variables only in secure deployment environments.

### Server Configuration

- **HTTP Server Defaults**:
  - Way's `New()` function sets safe HTTP server timeouts by default:
    - `ReadHeaderTimeout: 5s`
    - `ReadTimeout: 15s`
    - `WriteTimeout: 15s`
    - `IdleTimeout: 30s`
  - These defaults protect against slowloris and other timeout-related DoS attacks.
  - Custom timeouts can be set via `SetServer()` but should be reviewed carefully.

- **TLS/HTTPS**:
  - Way provides HTTP server primitives but does not enforce TLS.
  - Always run Way applications behind a reverse proxy with TLS termination (e.g., Traefik, Nginx) or configure TLS directly in `way.Server.TLSConfig`.

### Database Security

- **Connection Strings**:
  - Never hardcode database credentials in source code.
  - Use environment variables (e.g., `WAY_DB_USER`, `WAY_DB_PASSWORD`) or secure secret management.
  - Do not log connection strings or query parameters containing sensitive data.
  - Import only the driver adapter package your application needs, such as `github.com/swayedev/way/database/drivers/sqlite`.

- **SQL Injection**:
  - Always use parameterized queries. Way's `DB` helpers support parameterized queries via `Query()`, `QueryRow()`, and `Exec()`.
  - Avoid string concatenation to build SQL queries.

### Middleware and Error Handling

- **Error Responses**:
  - Do not expose stack traces or internal error details in production HTTP responses.
  - Use error middleware to safely render errors to clients.

- **Request Recovery**:
  - Implement middleware to recover from panics and log them safely without exposing to clients.

- **CORS and Security Headers**:
  - Implement CORS and security headers as middleware (e.g., `X-Frame-Options`, `X-Content-Type-Options`, `Content-Security-Policy`).
  - These are not included by default; add them as needed for your application.

### Logging

- **Sensitive Data**:
  - Do not log SQL query parameters, request headers, cookies, or authentication tokens.
  - Use structured logging to redact sensitive information.
  - Way's logging middleware logs only method, path, and duration by default.
  - Way's SQL helpers log operation status and errors, not raw query text or arguments.

### Outbound Requests

- `Context.ProxyMedia` uses Way's configured `HTTPClient` with a 15 second timeout by default.
- Use `SetHTTPClient()` to set shorter timeouts, custom transports, or network controls appropriate for your deployment.

### Dependencies

- **Minimal Dependencies**:
  - Way depends on `gorilla/mux`, `gorilla/securecookie`, `gorilla/sessions`, and `fcrypt`.
  - Keep these dependencies up to date and monitor for security advisories via `go list -u -m all` and tools like `govulncheck`.

### Deployment

- **Configuration**:
  - Run Way behind a reverse proxy with proper network segmentation.
  - Use environment variables for all configuration (database, keys, TLS, etc.).
  - Run Way with minimal privileges (non-root user).

- **Monitoring**:
  - Monitor application and access logs for anomalies.
  - Use rate limiting and DDoS protection at the infrastructure level.
  - Regularly audit database access logs.

## Security Updates

Way follows the Go community's security practices. Security patches for the current and previous minor versions will be backported as needed. Users are encouraged to upgrade to the latest stable version promptly.

## Known Limitations

- Way does not provide built-in rate limiting, authentication, or authorization. Implement these as middleware or at the infrastructure level.
- Way is designed for internal and moderate-traffic applications. For high-traffic or distributed systems, consider using more specialized frameworks.
- Session and cookie management are not hardened against all attack vectors (e.g., CSRF, session fixation). Implement additional middleware as needed.

## References

- [fcrypt Threat Model](https://github.com/swayedev/fcrypt/blob/main/THREAT_MODEL.md)
- [OWASP Top 10](https://owasp.org/Top10/)
- [Go Security Best Practices](https://golang.org/doc/effective_go#security)
