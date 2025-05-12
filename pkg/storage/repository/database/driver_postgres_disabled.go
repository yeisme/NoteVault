//go:build !postgres && !all

package database

import "fmt"

// PostgreSQL驱动未启用
func newPostgresDriverSafe() (DBDriver, error) {
	return nil, fmt.Errorf("PostgreSQL driver is not enabled in this build, please rebuild with -tags postgres or all")
}
