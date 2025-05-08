package model

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"phpbb-golang/internal/logger"
)

func OpenDb(ctx context.Context, tableName string) *sql.DB {
	// Refs:
	//   - https://go.dev/doc/tutorial/database-access
	//   - https://www.phpbb.com/demo/
	dbDSN := fmt.Sprintf("file:./model/db/%s.db?_foreign_keys=on", tableName)
	db, err := sql.Open("sqlite3", dbDSN)
	if err != nil {
		logger.Fatalf(ctx, "Error while opening %s table on database %s: %s", tableName, dbDSN, err)
	}
	return db
}

func DropDb(ctx context.Context, tableName string) error {
	// This is destructive!!
	db := OpenDb(ctx, tableName)
	defer db.Close()
	_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", Escape(tableName)))
	if err != nil {
		return fmt.Errorf("Error while dropping %s table: %s", tableName, err)
	}
	return nil
}

func Escape(sql string) string {
	// Escape the SQL data so that it is safe to use in query string. Please use prepared statement instead.
	// Ref: https://stackoverflow.com/questions/31647406/mysql-real-escape-string-equivalent-for-golang
	dest := make([]byte, 0, 2*len(sql))
	var escape byte
	for i := 0; i < len(sql); i++ {
		c := sql[i]
		escape = 0
		switch c {
		case 0: /* Must be escaped for 'mysql' */
			escape = '0'
			break
		case '\n': /* Must be escaped for logs */
			escape = 'n'
			break
		case '\r':
			escape = 'r'
			break
		case '\\':
			escape = '\\'
			break
		case '\'':
			escape = '\''
			break
		case '"': /* Better safe than sorry */
			escape = '"'
			break
		case '\032': //十进制26,八进制32,十六进制1a, /* This gives problems on Win32 */
			escape = 'Z'
		}
		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}
	return string(dest)
}
