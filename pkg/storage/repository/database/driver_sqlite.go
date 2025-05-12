//go:build sqlite || all

package database

// SQLite驱动已启用
func newSQLiteDriverSafe() (DBDriver, error) {
	return newSQLiteDriver(), nil
}
