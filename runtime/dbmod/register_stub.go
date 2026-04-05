//go:build !cgo

package mbdb

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const stubHint = "DB.* SQLite natives require CGO: set CGO_ENABLED=1, install gcc/MinGW and SQLite dev headers if needed, then rebuild"

func registerDBCommands(m *Module, reg runtime.Registrar) {
	_ = m
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			_ = args
			return value.Nil, fmt.Errorf("%s: %s", name, stubHint)
		}
	}
	keys := []string{
		"DB.OPEN", "DB.CLOSE", "DB.ISOPEN", "DB.EXEC", "DB.QUERY", "DB.QUERYJSON",
		"ROWS.NEXT", "ROWS.CLOSE", "ROWS.GETSTRING", "ROWS.GETINT", "ROWS.GETFLOAT",
		"DB.PREPARE", "DB.STMTCLOSE", "DB.STMTEXEC",
		"DB.BEGIN", "DB.COMMIT", "DB.ROLLBACK",
		"DB.LASTINSERTID", "DB.CHANGES",
	}
	for _, k := range keys {
		reg.Register(k, "db", stub(k))
	}
}
