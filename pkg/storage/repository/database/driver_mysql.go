//go:build mysql || all

package database

// MySQL驱动已启用
func newMySQLDriverSafe() (DBDriver, error) {
	return newMySQLDriver(), nil
}
