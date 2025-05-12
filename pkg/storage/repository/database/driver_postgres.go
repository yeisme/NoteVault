//go:build postgres || all

package database

// PostgreSQL驱动已启用
func newPostgresDriverSafe() (DBDriver, error) {
	return newPostgresDriver(), nil
}
