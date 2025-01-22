package sqlmanager

import (
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/ystyle/sqltemplate"
)

type SqlManager struct {
	sqlTemples []SqlTemple
	drivers    map[string]Driver
	funcs      sqltemplate.FuncMap
	tpl        *sqltemplate.Template
}

func New() *SqlManager {
	sm := &SqlManager{
		drivers: make(map[string]Driver),
		funcs:   sqltemplate.FuncMap{},
	}
	return sm
}

func (sm *SqlManager) Use(plugin Driver) {
	if _, ok := sm.drivers[plugin.DriverName()]; ok {
		log.Printf("sqlmanager - WARN: %s already used\n", plugin.DriverName())
	}
	sm.drivers[plugin.DriverName()] = plugin
}

func (sm *SqlManager) Load() {
	sm.tpl = sqltemplate.New("_root_")
	sm.tpl.Funcs(sm.funcs)
	var sqltpls []SqlTemple
	for _, driver := range sm.drivers {
		sqls, err := driver.Load()
		if err != nil {
			log.Printf("sqlmanager - ERROR: %s load failed: ", sqls)
			log.Panicln(err)
		}
		for _, sql := range sqls {
			d, has := sm.findTpl(sql.Name)
			if has {
				log.Printf("sqlmanager - WARN: %s Has duplicate sql: It will be cover [%s-%s] with [ %s ]\n", driver.DriverName(), sql.Name, strings.ReplaceAll(d.Sql, "\n", ""), strings.ReplaceAll(sql.Sql, "\n", ""))
			}
			sqltpls = append(sqltpls, sql)
			sm.tpl = sm.tpl.New(sql.Name)
			sm.tpl, err = sm.tpl.Parse(sql.Sql)
			if err != nil {
				panic(err)
			}
		}
		log.Printf("sqlmanager - INFO: %s loaded %d sqls.\n", driver.DriverName(), len(sqls))
	}
	sm.sqlTemples = sqltpls
}

func (sm *SqlManager) findTpl(name string) (*SqlTemple, bool) {
	for _, tpl := range sm.sqlTemples {
		if tpl.Name == name {
			return &tpl, true
		}
	}
	return nil, false
}

func (sm *SqlManager) RegisterFunc(funcs sqltemplate.FuncMap) {
	for k, v := range funcs {
		if temp, ok := sm.funcs[k]; ok {
			log.Printf("sqlmanager - WARN: %s Has duplicate func: It will be cover [%s] with [%s]\n", k, getFunctionName(temp), getFunctionName(v))
		}
		sm.funcs[k] = v
	}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
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
