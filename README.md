
# Logger Package

A simple, yet versatile logging library in Go, using `zap` for structured logging.

## Features

- Log messages with different levels: Debug, Info, Warn, Error, and Fatal.
- Configurable log output to various destinations (e.g., files, stdout).
- Structured logging with custom fields.
- Easily extendable to support different logging backends.

## Installation

To install the logger package, run:

```bash
go get github.com/ralonr/logger
```

## Usage

### Basic Usage

Here's an example of how to use the logger package to log messages to a file:

```go
package main

import (
    "os"
    "runtime"
    "github.com/ralonr/logger"
)

func main() {
    file, _ := os.OpenFile("example.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    var log logger.Logger = logger.NewZap(logger.Config{
        Level:  logger.DebugLevel,
        Output: file,
    })

    log.Debug("This is my debug", logger.Fields{
        "number of cpu": runtime.NumCPU(),
        "other info": 100,
    })
}
```

### Logging to Stdout

Here's an example of how to log messages to stdout:

```go
package main

import (
    "os"
    "runtime"
    "github.com/ralonr/logger"
)

func main() {
    var log logger.Logger = logger.NewZap(logger.Config{
        Level:  logger.DebugLevel,
        Output: os.Stdout,
    })

    log.Debug("This is my debug", logger.Fields{
        "number of cpu": runtime.NumCPU(),
        "other info": 100,
    })
}
```

### Custom Logger Implementation

You can create and use your own logger implementation by satisfying the `Logger` interface.

```go
package main

import (
    "fmt"
    "github.com/ralonr/logger"
)

type CustomLogger struct{}

func (c CustomLogger) Debug(msg string, fields logger.Fields) {
    fmt.Println("DEBUG:", msg, fields)
}

func (c CustomLogger) Info(msg string, fields logger.Fields) {
    fmt.Println("INFO:", msg, fields)
}

func (c CustomLogger) Warn(msg string, fields logger.Fields) {
    fmt.Println("WARN:", msg, fields)
}

func (c CustomLogger) Error(msg string, fields logger.Fields) {
    fmt.Println("ERROR:", msg, fields)
}

func (c CustomLogger) Fatal(msg string, fields logger.Fields) {
    fmt.Println("FATAL:", msg, fields)
}

func main() {
    var log logger.Logger = CustomLogger{}

    log.Debug("This is my debug", logger.Fields{
        "number of cpu": 4,
        "other info": 100,
    })
}
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
