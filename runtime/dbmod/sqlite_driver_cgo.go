//go:build cgo

package mbdb

import _ "github.com/mattn/go-sqlite3"

const sqliteDriverName = "sqlite3"
