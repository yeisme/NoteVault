//go:build !sqlite && !all

package database

import "fmt"

// SQLite驱动未启用
func newSQLiteDriverSafe() (DBDriver, error) {
	return nil, fmt.Errorf("SQLite driver is not enabled in this build, please rebuild with -tags sqlite or all")
}
