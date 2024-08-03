### Change log

All notable changes to this project will be documented in this file.

#### [Unreleased]

##### Added
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

#### [0.3.1] - 2023-12-31
##### Added
- **Initial Release (Pre Changelog)**
  - Basic driver checks and DSN construction for various databases.
  - Connection setup using `pgx` without connection pooling.
  - Simple query execution functions for `pgx`.
  - Basic connection setup using `database/sql` without connection pooling.
  - Simple query execution functions for `database/sql`.
