package model

import (
	"context"
	"database/sql"
	"fmt"
	"phpbb-golang/internal/logger"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDb(ctx context.Context, tableName string) *sql.DB {
	// Refs:
	//   - https://go.dev/doc/tutorial/database-access
	//   - https://www.phpbb.com/demo/
	dbDSN := fmt.Sprintf("file:./model/db/%s.db?_foreign_keys=on", tableName)
	db, err := sql.Open("sqlite3", dbDSN)
	if err != nil {
		logger.Fatalf(ctx, "Error while opening Database DSN %s: %s", dbDSN, err)
	}
	return db
}
