package sqlmanager

import (
	"embed"
	"fmt"
	"testing"
)

//go:embed test-sql
var Assets embed.FS

func TestEmbedMarkdownDriver(t *testing.T) {
	sm := New()
	sm.Use(NewMarkdownDriverWithEmbedDir(Assets, "test-sql"))
	// sm.Use(sqlmanager.NewMarkdownDriverWithEmbed(Assets)) // default dir is sql
	sm.Load()
	sql, args, _ := sm.RenderTPL("test/GetStudentByID2", 1)
	fmt.Println(sql)
	fmt.Println(args)
}
