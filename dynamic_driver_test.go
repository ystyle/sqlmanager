package sqlmanager

import (
	"fmt"
	"strings"
	"testing"
	"text/template"
)

func TestDynamicLoad(t *testing.T) {
	sm := New()
	store := NewDynamicDriver()
	store.Register("rest1", `select * from table where id = {{. }} or {{ block "rest" . }} {{ end }}`)
	store.Register("rest", `select * from table where id = "{{ upper . }}"`)
	sm.Use(store)
	sm.RegisterFunc(template.FuncMap{
		"upper": func(v string) string {
			return strings.ToUpper(v)
		},
	})
	sm.Load()
	sql, args, err := sm.RenderTPL("rest1", "test")
	if err != nil {
		panic(err)
	}
	fmt.Println(sql)
	fmt.Println(args)
}
