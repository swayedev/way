### Changelog for Way Framework

All notable changes to this project will be documented in this file.

#### [Unreleased]

##### Added
- **Config Struct**: Introduced a `Config` struct to hold all configuration parameters, enabling a more structured and flexible configuration management system.
- **Default Config**: Implemented `defaultConfig` function to generate default configurations if not provided, ensuring the framework is always initialized with valid settings.
- **Random Key Generation**: Added functionality to generate random keys for store and cookie encryption when not provided, enhancing security by default.
- **Way Context**: Added `Bind` and `MultipartForm` to enhance the way frameworks features similar to `echo framework`
- **Centralized Configuration**
  - Introduced `DriverConfig` struct for managing configuration of different database drivers.
  - Introduced the `DBConfig` interface to standardize configuration handling across different database drivers.
  - Introduced `PGXConfig` struct for centralized configuration and optional pooling, implementing the `DBConfig` interface for PostgreSQL databases.
  - Introduced `SQLConfig` struct for centralized configuration and optional pooling, implementing the `DBConfig` interface for SQL databases.

- **Custom Error Handling**
  - Implemented `DBError` type for consistent and detailed error handling across all modules.
  - Added operation constants for better error context (e.g., `OpConfigParse`, `OpDriverCheck`, `OpDSNSet`, `OpPing`, `OpConnect`, `OpExec`, `OpQuery`, `OpClose`, `OpPoolConnect`, `OpTxBegin`, `OpTxCommit`, `OpTxRollback`).

- **Optional Pooling**
  - Added support for optional connection pooling in `pgx.go`.
  - Added support for optional connection pooling in `sql.go`.

- **Transaction Management**
  - Added transaction management functions in `pgx.go` (`PGXBeginTx`, `PGXCommitTx`, `PGXRollbackTx`).
  - Added transaction management functions in `sql.go` (`SQLBeginTx`, `SQLCommitTx`, `SQLRollbackTx`).

- **Improved Documentation**
  - Added comments for better function documentation across all modules.

##### Changed
- **New Function**: Updated the `New` function to accept an optional `Config` struct. If no config is provided, the function generates a default configuration and random keys for encryption.
- **Logging**: Improved logging to include the generation of random encryption keys, while ensuring keys are not logged directly for security reasons.
- **Database Initialization**: Refactored `InitDBFromConfig` to use a helper function `initDBConnection`, reducing code duplication and improving readability.
- **Middleware Adaptation**: Updated `adaptMiddleware` function to better handle the adaptation of middleware functions to the `mux.MiddlewareFunc` type.
- **Session Management**: Improved session initialization by checking if default session values should be used and setting them accordingly.

- **Database Driver Checking**
  - Improved `CheckDriver` function to return detailed errors using `DBError`.

- **DSN Construction**
  - Consolidated and improved DSN construction functions with detailed error handling.
  - Enhanced DSN construction logic to use `DriverConfig` for centralized configuration.

- **Query Execution**
  - Consolidated query execution logic in `pgx.go` for both pooled and non-pooled connections.
  - Consolidated query execution logic in `sql.go`.

- **Logging and Error Messages**
  - Improved logging and error messages for better debugging and maintainability across all modules.

### Deprecated
- No deprecations in this version.

### Removed
- No removals in this version.

### Fixed
- **Typo in ASCII Art**: Corrected a typo in the ASCII art displayed when the server starts.

#### [0.3.1] - 2023-12-31
##### Added
- **Initial Release (Pre Changelog)**
  - **Way Struct**: Core structure encapsulating the router, server, logger, sessions, and database connections.
  - **HandlerFunc and MiddlewareFunc Types**: Defined custom handler and middleware function types for streamlined request handling.
  - **New Function**: Initializes a new `Way` instance, sets up session defaults, and configures the router and logger.
  - **Logger Management**: Functions to set and get the logger.
  - **Router and Server Management**: Functions to set the router and server instances.
  - **Listener Management**: Functions to set the network listener.
  - **Database Management**: Functions to set and initialize the database connection from environment variables.
  - **Session Management**: Functions to set and manage sessions.
  - **Middleware Management**: `Use` function to add middleware to the middleware stack.
  - **Request Handling**: Functions to register routes (`HandleFunc`, `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `OPTIONS`, `HEAD`).
  - **Server Control**: Functions to start, close, and gracefully shutdown the server.
  - **Utility Functions**: Functions to retrieve environment variables and create default loggers.
  - **Session Defaults**: Functions to check and set default session values.
  - Basic driver checks and DSN construction for various databases.
  - Connection setup using `pgx` without connection pooling.
  - Simple query execution functions for `pgx`.
  - Basic connection setup using `database/sql` without connection pooling.
  - Simple query execution functions for `database/sql`.
  - Added constants for environment variable keys, making it easier to manage and update environment variables.
