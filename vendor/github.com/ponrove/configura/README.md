# Configura - Go Shared Configuration Management

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/e21f682950c3464c9f241806ae5d892d)](https://app.codacy.com/gh/ponrove/configura/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/e21f682950c3464c9f241806ae5d892d)](https://app.codacy.com/gh/ponrove/configura/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)
[![GoDoc](https://godoc.org/github.com/ponrove/configura?status.svg)](https://godoc.org/github.com/ponrove/configura)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`configura` is a Go package designed to simplify application configuration management. It provides a type-safe way to define, load, and access configuration variables from environment variables, and to share these configurations across modules or subpackages in your Go application.

## Description

Managing configuration variables, especially across different environments (development, staging, production), can be error-prone. `configura` addresses this by:

- **Type Safety:** Defining configuration variables with specific Go types (e.g., `string`, `int`, `bool`). This helps catch errors at compile-time or during setup rather than runtime.
- **Centralized Definition:** Encouraging the definition of all expected configuration variables.
- **Environment Variable Loading:** Easily loading values from environment variables with fallbacks for missing ones.
- **Validation:** Allowing components or subpackages to declare their required configuration keys and verify their presence.

## Usage

### 1. Define Your Configuration Variables

It's good practice to define your configuration `Variable` constants in a dedicated package or a specific part of your application.

```go
package config

import "github.com/ponrove/configura"

// Define your application's configuration variables
const (
	DATABASE_URL     configura.Variable[string] = "DATABASE_URL"
	PORT             configura.Variable[int]    = "PORT"
	API_KEY          configura.Variable[string] = "API_KEY"
	ENABLE_FEATURE_X configura.Variable[bool]   = "ENABLE_FEATURE_X"
	TIMEOUT_SECONDS  configura.Variable[int64]  = "TIMEOUT_SECONDS"
)
```

### 2. Initialize and Load Configuration

In your application's main setup (e.g., `main.go`), you'll initialize a `ConfigImpl` and load the environment variables.

```go
package main

import (
	"os"

	"github.com/ponrove/configura"
	"github.com/ponrove/configura/_example/config"
	"github.com/ponrove/configura/_example/subpackage"
)

func main() {
	// --- Simulate setting environment variables (for example purposes) ---
	// In a real scenario, these would be set in your shell, Dockerfile, K8s manifest, etc.
	os.Setenv(string(config.DATABASE_URL), "postgres://user:pass@host:port/dbname")
	os.Setenv(string(config.PORT), "8080")
	os.Setenv(string(config.API_KEY), "supersecretapikey")
	os.Setenv(string(config.ENABLE_FEATURE_X), "true")
	os.Setenv(string(subpackage.SUBPACKAGE_DEFINED_CONFIG), "some_value")
	// TIMEOUT_SECONDS is not set, so its fallback will be used.

	// --- Initialize Configura ---
	cfg := configura.NewConfigImpl()

	// Load environment variables with fallbacks
	configura.LoadEnvironment(cfg, config.DATABASE_URL, "postgres://fallback_user:fallback_pass@localhost:5432/fallback_db")
	configura.LoadEnvironment(cfg, config.PORT, 3000)  // Fallback port 3000
	configura.LoadEnvironment(cfg, config.API_KEY, "") // Fallback empty string if not set
	configura.LoadEnvironment(cfg, config.ENABLE_FEATURE_X, false)
	configura.LoadEnvironment(cfg, config.TIMEOUT_SECONDS, int64(30)) // Fallback 30 seconds
	configura.LoadEnvironment(cfg, subpackage.SUBPACKAGE_DEFINED_CONFIG, "default_value")

	// Set the configuration by yourself
	cfg.RegInt64[config.TIMEOUT_SECONDS] = int64(25)

	err := subpackage.Initialize(cfg)
	if err != nil {
		panic(err) // Handle error appropriately in your application
	}
}
```

### 3. Subpackage Configuration Validation

A subpackage (e.g., `subpackage`) can ensure that all configuration variables it depends on are present in the `configura.Config` instance it receives.

```go
// File: _example/subpackage.go
package subpackage

import (
	"fmt"

	"github.com/ponrove/configura"
	"github.com/ponrove/configura/_example/config"
)

const (
	SUBPACKAGE_DEFINED_CONFIG configura.Variable[string] = "SUBPACKAGE_DEFINED_CONFIG"
)

// RequiredUserServiceKeys lists the configuration variables this service needs.
var RequiredUserServiceKeys = []any{
	SUBPACKAGE_DEFINED_CONFIG,
	config.DATABASE_URL,
	config.API_KEY,
}

// Initialize sets up the user service with the given configuration.
// It validates that all required configuration keys are registered.
func Initialize(cfg configura.Config) error {
	// Validate that the config instance has all the keys our service needs
	if err := cfg.ConfigurationKeysRegistered(RequiredUserServiceKeys...); err != nil {
		return fmt.Errorf("user service configuration validation failed: %w", err)
	}

	// Access the validated configuration
	dbURL := cfg.String(config.DATABASE_URL)
	apiKey := cfg.String(config.API_KEY)
	definedConfig := cfg.String(SUBPACKAGE_DEFINED_CONFIG)

	fmt.Printf("UserService: Initializing with DB URL: %s and API Key (present: %t), and has subpackage defined key (present: %s)\n", dbURL, apiKey != "", definedConfig)
	// ... further initialization logic for the user service ...

	return nil
}
```

### How `ConfigurationKeysRegistered` Works

The `ConfigurationKeysRegistered` method iterates through the provided keys. If any key is not found in the `ConfigImpl`'s internal maps (meaning `LoadEnvironment` was not called for it, or it wasn't otherwise set), it returns a `ErrMissingVariable`. This error contains a list of all the missing keys.

This allows for robust startup checks, ensuring your application components have the configuration they need before they start running.

## Contributing

Contributions are welcome! Please feel free to open a pull request with any improvements, bug fixes, or new features.

1.  Fork the repository.
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.

Please ensure your code is well-tested and follows Go best practices.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
