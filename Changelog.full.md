## Change Log

#### `database.go`

**Initial Version:**
- Basic driver checks and DSN construction for various databases.
- Limited error handling with basic error messages.
- Lack of centralized configuration and context management.

**Current Version:**
- Introduced `DriverConfig` struct to hold configuration for database drivers.
- Added `DBError` type for custom error handling with constants for operations.
- Implemented centralized error messages.
- Improved function documentation with comments.
- Enhanced DSN construction with more detailed error handling.
- Consolidated driver checks and DSN setting logic.

**Changes:**
```markdown
- Added `DriverConfig` struct for centralized configuration.
- Introduced `DBError` type for improved error handling.
- Added operation constants (e.g., `OpConfigParse`, `OpDriverCheck`, `OpDSNSet`).
- Consolidated and improved DSN construction functions with better error messages.
- Updated driver check function to return detailed errors.
- Improved documentation and code comments.
```

#### `pgx.go`

**Initial Version:**
- Basic connection setup using `pgx` without connection pooling.
- Simple query execution functions.
- Limited error handling.

**Current Version:**
- Introduced `PGXConfig` struct to manage connection settings.
- Added support for optional connection pooling.
- Implemented `DBError` for custom error handling.
- Consolidated query execution logic for pooled and non-pooled connections.
- Added functions for transaction management.
- Improved logging and error messages.

**Changes:**
```markdown
- Added `PGXConfig` struct for centralized configuration and pooling options.
- Introduced support for optional connection pooling.
- Enhanced error handling using `DBError`.
- Consolidated query execution logic for both pooled and non-pooled connections.
- Added functions for transaction management (`PGXBeginTx`, `PGXCommitTx`, `PGXRollbackTx`).
- Improved logging and error messages.
- Updated function documentation with comments.
```

#### `sql.go`

**Initial Version:**
- Basic connection setup using `database/sql` without connection pooling.
- Simple query execution functions.
- Limited error handling.

**Current Version:**
- Introduced `SQLConfig` struct to manage connection settings.
- Added support for optional connection pooling.
- Implemented `DBError` for custom error handling.
- Consolidated query execution logic.
- Added functions for transaction management.
- Improved logging and error messages.
- Enhanced function documentation with comments.

**Changes:**
```markdown
- Added `SQLConfig` struct for centralized configuration and pooling options.
- Introduced support for optional connection pooling.
- Enhanced error handling using `DBError`.
- Consolidated query execution logic.
- Added functions for transaction management (`SQLBeginTx`, `SQLCommitTx`, `SQLRollbackTx`).
- Improved logging and error messages.
- Updated function documentation with comments.
```
