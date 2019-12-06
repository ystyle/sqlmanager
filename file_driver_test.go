package sqlmanager

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	sm := New()
	sm.Use(NewMarkdownDriverWithDir("./test-sql"))
	sm.Load()
	sql, err := sm.RenderTPL("GetStudentByID2", 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(sql)
}
