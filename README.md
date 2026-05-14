# Way Framework

![Way Mascot](./way_mascot.png)

[![Go Reference](https://pkg.go.dev/badge/github.com/swayedev/way.svg)](https://pkg.go.dev/github.com/swayedev/way)

_Version: 1.0.0-rc1_

**Status:** Production-ready release candidate. Way v1.0.0-rc1 is the first stable API release. See [API_FREEZE.md](API_FREEZE.md) for the frozen public API surface.

## Overview

Way is a lightweight, Go-based web framework that integrates with the [Gorilla Mux](https://github.com/gorilla/mux) router and provides simplified mechanisms for handling HTTP requests inspired by the [echo framework](https://echo.labstack.com) while adding database operations, session management, and cryptographic helpers.

Way is designed for building internal and moderate-traffic web applications. It pairs well with the [fcrypt](https://github.com/swayedev/fcrypt) encryption library for secure data protection.

## Documentation

- [API Freeze](API_FREEZE.md) – frozen public APIs and stability guarantees
- [Security Policy](SECURITY.md) – security considerations and hardening guidance  
- [Migration Guide](MIGRATION.md) – upgrade from pre-1.0 versions
- [Changelog](CHANGELOG.md) – release history

## Features

- Custom context for HTTP handlers with response helpers (JSON, XML, HTML, String, Data, Images)
- Simplified route declaration (GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD)
- Integrated SQL database operations with support for MySQL, PostgreSQL (pgx), SQLite, SQL Server, and Oracle
- Session and cookie management with Gorilla Sessions
- Graceful shutdown and startup management with safe HTTP server timeouts
- Encryption and hashing via [fcrypt](https://github.com/swayedev/fcrypt) integration
- Request logging middleware with method, path, and duration tracking

## Getting Started

### Installation
To start using the Way framework, install it using `go get`:
```bash
go get -u github.com/swayedev/way@v1.0.0-rc1
```

### Basic Usage
Here's a simple example to get you started:

```go
package main

import (
    "github.com/swayedev/way"
)

func main() {
    w := way.New()

    w.GET("/", func(ctx *way.Context) {
        ctx.Response.Write([]byte("Hello, World!"))
    })

    w.Start(":8080")
}
```

Or with a handler function

```go
package main

import "github.com/swayedev/way"

func main() {
	w := way.New()
	w.GET("/", helloHandler)

	w.Start(":8080")
}

func helloHandler(c *way.Context) {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.Write([]byte("Hello World"))
}

```

## Routing
Way simplifies route handling with predefined methods for standard HTTP verbs:

```go
w.GET("/path", yourGetHandler)
w.POST("/path", yourPostHandler)
// ... and so on for PUT, DELETE, PATCH, OPTIONS, HEAD
```

## Database Operations

### Passing an existing sql db

This section explains how to pass an existing SQL database to the application. It provides instructions and guidelines on how to configure the application to use an existing database instead of creating a new one.

```go


```

### Opening a Connection
```go
err := w.SqlOpen()
if err != nil {
    // Handle error
}
```

### Executing Queries
```go
result, err := w.SqlExec(context.Background(), "your SQL query here", args...)
if err != nil {
    // Handle error
}
```

### Querying Data
```go
rows, err := w.SqlQuery(context.Background(), "your SQL query here", args...)
if err != nil {
    // Handle error
}
// Remember to close rows
```

## Graceful Shutdown
To gracefully shut down your server:

```go
err := w.Shutdown(context.Background())
if err != nil {
    // Handle error
}
```

## Contributing
Contributions to the Way framework are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License
This project is licensed under the [MIT License](LICENSE).
