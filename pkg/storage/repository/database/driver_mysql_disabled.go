//go:build !mysql && !all

package database

import "fmt"

// MySQL驱动未启用
func newMySQLDriverSafe() (DBDriver, error) {
	return nil, fmt.Errorf("MySQL driver is not enabled in this build, please rebuild with -tags mysql or all")
}
