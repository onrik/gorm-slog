# gorm-slog
[![Go Report Card](https://goreportcard.com/badge/github.com/onrik/gorm-slog)](https://goreportcard.com/report/github.com/onrik/gorm-slog)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/onrik/gorm-slog)](https://pkg.go.dev/github.com/onrik/gorm-slog)

Logger for gorm based on log/slog.

```golang
package main

import (
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
    "github.com/onrik/gorm-slog"
)

func main() {
    db, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{
        Logger: gormslog.New(nil),
    })
    if err != nil {
        panic("failed to connect database")
    }
}

```
