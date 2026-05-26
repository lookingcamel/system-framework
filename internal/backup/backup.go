package backup

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

func BackupDatabase(db *sql.DB, dbType string, backupDir string) (string, error) {
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("backup_%s.sql", timestamp))

	switch dbType {
	case "mysql":
		return backupFile, backupMySQL(db, backupFile)
	case "postgres", "postgresql":
		return backupFile, backupPostgreSQL(db, backupFile)
	case "sqlite":
		return backupFile, backupSQLite(db, backupFile)
	default:
		return backupFile, backupSQLite(db, backupFile)
	}
}

func backupMySQL(db *sql.DB, backupFile string) error {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return err
		}
		tables = append(tables, table)
	}

	file, err := os.Create(backupFile)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, table := range tables {
		_, err := file.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS `%s`;\n", table))
		if err != nil {
			return err
		}

		createTable, err := db.Query(fmt.Sprintf("SHOW CREATE TABLE `%s`", table))
		if err != nil {
			return err
		}
		if createTable.Next() {
			var t, createStmt string
			if err := createTable.Scan(&t, &createStmt); err != nil {
				return err
			}
			_, err := file.WriteString(createStmt + ";\n\n")
			if err != nil {
				return err
			}
		}
		createTable.Close()

		dataRows, err := db.Query(fmt.Sprintf("SELECT * FROM `%s`", table))
		if err != nil {
			return err
		}
		cols, err := dataRows.Columns()
		if err != nil {
			return err
		}

		var values []string
		for dataRows.Next() {
			rowData := make([]interface{}, len(cols))
			rowPointers := make([]interface{}, len(cols))
			for i := range rowData {
				rowPointers[i] = &rowData[i]
			}
			if err := dataRows.Scan(rowPointers...); err != nil {
				return err
			}

			var rowValues []string
			for _, val := range rowData {
				if val == nil {
					rowValues = append(rowValues, "NULL")
				} else {
					rowValues = append(rowValues, fmt.Sprintf("'%v'", val))
				}
			}
			values = append(values, fmt.Sprintf("(%s)", joinValues(rowValues)))
		}
		dataRows.Close()

		if len(values) > 0 {
			_, err := file.WriteString(fmt.Sprintf("INSERT INTO `%s` VALUES %s;\n\n", table, joinValues(values)))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func backupPostgreSQL(db *sql.DB, backupFile string) error {
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		return err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return err
		}
		tables = append(tables, table)
	}

	file, err := os.Create(backupFile)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, table := range tables {
		createTable, err := db.Query(fmt.Sprintf("SELECT pg_get_tabledef('%s')", table))
		if err != nil {
			return err
		}
		if createTable.Next() {
			var createStmt string
			if err := createTable.Scan(&createStmt); err != nil {
				return err
			}
			_, err := file.WriteString(createStmt + ";\n\n")
			if err != nil {
				return err
			}
		}
		createTable.Close()

		dataRows, err := db.Query(fmt.Sprintf("SELECT * FROM \"%s\"", table))
		if err != nil {
			return err
		}
		cols, err := dataRows.Columns()
		if err != nil {
			return err
		}

		var values []string
		for dataRows.Next() {
			rowData := make([]interface{}, len(cols))
			rowPointers := make([]interface{}, len(cols))
			for i := range rowData {
				rowPointers[i] = &rowData[i]
			}
			if err := dataRows.Scan(rowPointers...); err != nil {
				return err
			}

			var rowValues []string
			for _, val := range rowData {
				if val == nil {
					rowValues = append(rowValues, "NULL")
				} else {
					rowValues = append(rowValues, fmt.Sprintf("'%v'", val))
				}
			}
			values = append(values, fmt.Sprintf("(%s)", joinValues(rowValues)))
		}
		dataRows.Close()

		if len(values) > 0 {
			_, err := file.WriteString(fmt.Sprintf("INSERT INTO \"%s\" VALUES %s;\n\n", table, joinValues(values)))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func backupSQLite(db *sql.DB, backupFile string) error {
	file, err := os.Create(backupFile)
	if err != nil {
		return err
	}
	defer file.Close()

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		return err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return err
		}
		if table != "sqlite_sequence" {
			tables = append(tables, table)
		}
	}

	for _, table := range tables {
		createTable, err := db.Query(fmt.Sprintf("SELECT sql FROM sqlite_master WHERE name='%s'", table))
		if err != nil {
			return err
		}
		if createTable.Next() {
			var createStmt string
			if err := createTable.Scan(&createStmt); err != nil {
				return err
			}
			_, err := file.WriteString(createStmt + ";\n\n")
			if err != nil {
				return err
			}
		}
		createTable.Close()

		dataRows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", table))
		if err != nil {
			return err
		}
		cols, err := dataRows.Columns()
		if err != nil {
			return err
		}

		var values []string
		for dataRows.Next() {
			rowData := make([]interface{}, len(cols))
			rowPointers := make([]interface{}, len(cols))
			for i := range rowData {
				rowPointers[i] = &rowData[i]
			}
			if err := dataRows.Scan(rowPointers...); err != nil {
				return err
			}

			var rowValues []string
			for _, val := range rowData {
				if val == nil {
					rowValues = append(rowValues, "NULL")
				} else {
					rowValues = append(rowValues, fmt.Sprintf("'%v'", val))
				}
			}
			values = append(values, fmt.Sprintf("(%s)", joinValues(rowValues)))
		}
		dataRows.Close()

		if len(values) > 0 {
			_, err := file.WriteString(fmt.Sprintf("INSERT INTO %s VALUES %s;\n\n", table, joinValues(values)))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func joinValues(values []string) string {
	result := ""
	for i, v := range values {
		if i > 0 {
			result += ", "
		}
		result += v
	}
	return result
}
