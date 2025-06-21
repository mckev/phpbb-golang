package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/mattn/go-sqlite3"

	"phpbb-golang/internal/logger"
)

const (
	DB_ERROR_NO_RESULT         = "[DB_ERROR_NO_RESULT]"
	DB_ERROR_UNIQUE_CONSTRAINT = "[DB_ERROR_UNIQUE_CONSTRAINT]"
)

func OpenDb(ctx context.Context, tableName string) *sql.DB {
	// Refs:
	//   - https://go.dev/doc/tutorial/database-access
	//   - https://www.phpbb.com/demo/
	dbDSN := "file:./model/db/main.db?_foreign_keys=on"
	db, err := sql.Open("sqlite3", dbDSN)
	if err != nil {
		logger.Fatalf(ctx, "Error while opening table '%s' on database '%s': %s", tableName, dbDSN, err)
	}
	return db
}

func DropDb(ctx context.Context, tableName string) error {
	// WARNING: This is destructive!!
	db := OpenDb(ctx, tableName)
	defer db.Close()
	_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", SqlEscape(tableName)))
	if err != nil {
		return fmt.Errorf("Error while dropping table '%s': %s", tableName, err)
	}
	return nil
}

func SqlEscape(sql string) string {
	// Escape the SQL data so that it is safe to use in query string. Please use prepared statement instead!
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

func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	// SQLite (mattn/go-sqlite3)
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
	}

	// PostgreSQL (lib/pq)
	// var pqErr *pq.Error
	// if errors.As(err, &pqErr) {
	//     return pqErr.Code == "23505" // unique_violation
	// }

	// MySQL (go-sql-driver/mysql)
	// var mysqlErr *mysql.MySQLError
	// if errors.As(err, &mysqlErr) {
	//     return mysqlErr.Number == 1062 // ER_DUP_ENTRY
	// }

	// Fallback: parse error message (least reliable)
	// SQLite: UNIQUE constraint failed: users.user_name
	// PostgreSQL: pq: duplicate key value violates unique constraint "users_user_name_key"
	// MySQL: Error 1062: Duplicate entry 'alice' for key 'users.user_name'
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE constraint failed") || // SQLite
		strings.Contains(msg, "duplicate key value violates unique constraint") || // PostgreSQL
		strings.Contains(msg, "Error 1062: Duplicate entry") // MySQL
}
