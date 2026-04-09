//go:build !cgo && modernc_sqlite

package mbdb

import _ "modernc.org/sqlite"

const sqliteDriverName = "sqlite"
