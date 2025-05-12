package database

import "fmt"

func NewDatabaseDriver(driverName string) (DBDriver, error) {
	switch driverName {
	case "mysql":
		return newMySQLDriverSafe()
	case "postgres":
		return newPostgresDriverSafe()
	case "sqlite", "sqlite3":
		return newSQLiteDriverSafe()
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driverName)
	}
}
