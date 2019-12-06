package sqlmanager

import (
	"fmt"
	"testing"
)

func TestDynamicLoad(t *testing.T) {
	sm := New()
	store := NewDynamicDriver()
	store.Register("rest", `select * from table where id = {{.}}`)
	sm.Use(store)
	sm.Load()
	sql, err := sm.RenderTPL("rest", 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(sql)
}
