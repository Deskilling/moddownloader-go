package launcher

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/charmbracelet/log"
)

func DumpDatabase(path string) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal("open failed", "err", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("ping failed", "err", err)
	}
	tables, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'`)
	if err != nil {
		log.Fatal("failed to list tables", "err", err)
	}
	defer tables.Close()

	for tables.Next() {
		var table string
		if err := tables.Scan(&table); err != nil {
			log.Error("scan table failed", "err", err)
			continue
		}

		fmt.Println("=== Tabelle:", table, "===")

		cols, err := db.Query(fmt.Sprintf(`PRAGMA table_info(%s)`, table))
		if err != nil {
			log.Error("failed to get columns", "table", table, "err", err)
			continue
		}

		var colNames []string
		for cols.Next() {
			var cid int
			var name, ctype string
			var notnull, pk int
			var dflt any

			cols.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk)
			colNames = append(colNames, name)
		}
		cols.Close()

		query := fmt.Sprintf(`SELECT * FROM %s`, table)
		rows, err := db.Query(query)
		if err != nil {
			log.Error("failed query", "table", table, "err", err)
			continue
		}

		for rows.Next() {
			values := make([]any, len(colNames))
			valuePtrs := make([]any, len(colNames))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				log.Error("scan row failed", "table", table, "err", err)
				continue
			}

			rowMap := make(map[string]any)
			for i, col := range colNames {
				rowMap[col] = values[i]
			}

			fmt.Println(rowMap)
		}

		rows.Close()
		fmt.Println()
	}
}
