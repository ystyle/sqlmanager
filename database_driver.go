package sqlmanager

import (
	"database/sql"
	"fmt"
)

const sqlRaw = "select name, description, `sql` from %s where deleted_at is null"

type DatabaseDriver struct {
	db        *sql.DB
	tablename string
	sqls      []SqlTemple
}

func NewDatabaseDriver(db *sql.DB, tablename string) *DatabaseDriver {
	return &DatabaseDriver{db: db, tablename: tablename}
}

func (dbd *DatabaseDriver) DriverName() string {
	return "database"
}

func (dbd *DatabaseDriver) Load() ([]SqlTemple, error) {
	rows, err := dbd.db.Query(fmt.Sprintf(sqlRaw, dbd.tablename))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var sql SqlTemple
		if err := rows.Scan(&sql.Name, &sql.Description, &sql.Sql); err != nil {
			return nil, err
		}
		dbd.sqls = append(dbd.sqls, sql)
	}
	return dbd.sqls, nil
}
