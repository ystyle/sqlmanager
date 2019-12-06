package sqlmanager

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type sqlManager struct {
	sqlTemples []SqlTemple
	drivers    map[string]Driver
}

func New() *sqlManager {
	sm := &sqlManager{
		drivers: make(map[string]Driver),
	}
	return sm
}

func (sm *sqlManager) Use(plugin Driver) {
	sm.drivers[plugin.DriverName()] = plugin
}

func (sm *sqlManager) Load() {
	for _, driver := range sm.drivers {
		sqls, err := driver.Load()
		if err != nil {
			log.Printf("sqlmanager - ERROR: %s load failed: ", sqls)
			log.Panicln(err)
		}
		for _, sql := range sqls {
			d, err := sm.findTpl(sql.Name)
			if err != nil {
				sm.sqlTemples = append(sm.sqlTemples, sql)
			} else {
				log.Printf("sqlmanager - WARN: %s Has duplicate sql: %s [ %s ]", sql.Name, d.Name, strings.ReplaceAll(d.Sql, "\n", ""))
				log.Printf("sqlmanager - WARN: It will be covered")
			}
		}

		log.Printf("sqlmanager - INFO: %s loaded %d sqls.\n", driver.DriverName(), len(sqls))
	}
}

func (sm *sqlManager) findTpl(name string) (*SqlTemple, error) {
	for _, tpl := range sm.sqlTemples {
		if tpl.Name == name {
			return &tpl, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("sqlmanager - template: %s no found", name))
}

type Driver interface {
	Load() ([]SqlTemple, error)
	DriverName() string
}

type SqlTemple struct {
	Name        string // 名称
	Description string // 描述
	Sql         string // sql
}