package sqlmanager

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"log"
	"os"
	"path"
	"path/filepath"
)

type MarkdownDriver struct {
	dir string
}

type item struct {
	Type    string
	Content string
}

func NewMarkdownDriver() *MarkdownDriver {
	return &MarkdownDriver{
		dir: "./sql",
	}
}

func NewMarkdownDriverWithDir(dir string) *MarkdownDriver {
	return &MarkdownDriver{
		dir: dir,
	}
}
func (mdd *MarkdownDriver) DriverName() string {
	return "Markdown"
}

func (mdd *MarkdownDriver) Load() ([]SqlTemple, error) {
	var sqls []SqlTemple
	err := filepath.Walk(mdd.dir, func(subpath string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		ext := path.Ext(subpath)
		if ext == ".md" || ext == ".markdown" {
			s, err := parseMarkdown(subpath)
			if err != nil {
				return err
			}
			if len(s) != 0 {
				sqls = append(sqls, s...)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return sqls, nil
}

func parseMarkdown(filename string) ([]SqlTemple, error) {
	var sqls []SqlTemple
	buf, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("sqlmanager - ERROR: %s loading failed...\n", filename)
		return nil, err
	}
	if bytes.ContainsRune(buf, '\r') {
		buf = markdown.NormalizeNewlines(buf)
	}
	psr := parser.New()
	node := markdown.Parse(bytes.ReplaceAll(buf, []byte("\r"), nil), psr)
	list := getAll(node)
	i := 0
	for {
		// 1. text, code
		// 2. text, text, code
		if i >= len(list) {
			break
		}
		var tpl SqlTemple
		if list[i].Type == "text" && list[i+1].Type == "code" {
			tpl.Name = list[i].Content
			tpl.Sql = list[i+1].Content
			sqls = append(sqls, tpl)
			i += 2
		} else if list[i].Type == "text" && list[i+1].Type == "text" && list[i+2].Type == "code" {
			tpl.Name = list[i].Content
			tpl.Description = list[i+1].Content
			tpl.Sql = list[i+2].Content
			sqls = append(sqls, tpl)
			i += 3
		} else {
			return nil, errors.New(fmt.Sprintf("ERROR: parse markdown failed: %s", filename))
		}
	}
	return sqls, nil
}
