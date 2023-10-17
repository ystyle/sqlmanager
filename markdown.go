package sqlmanager

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

func parseMarkdown(buf []byte, prefix string) ([]SqlTemple, error) {
	var sqls []SqlTemple
	if bytes.ContainsRune(buf, '\r') {
		buf = markdown.NormalizeNewlines(buf)
	}
	psr := parser.New()
	node := markdown.Parse(buf, psr)
	list := getAll(node)
	i := 0
	for {
		// 1. text, code
		// 2. text, text, code
		if i >= len(list) {
			break
		}
		name := filepath.Join(prefix, list[i].Content)
		name = strings.ReplaceAll(name, "\\", "/")
		var tpl SqlTemple
		if list[i].Type == "text" && list[i+1].Type == "code" {
			tpl.Name = name
			tpl.Sql = list[i+1].Content
			sqls = append(sqls, tpl)
			i += 2
		} else if list[i].Type == "text" && list[i+1].Type == "text" && list[i+2].Type == "code" {
			tpl.Name = name
			tpl.Description = list[i+1].Content
			tpl.Sql = list[i+2].Content
			sqls = append(sqls, tpl)
			i += 3
		} else {
			return nil, errors.New("parse markdown failed")
		}
	}
	return sqls, nil
}

func getAll(node ast.Node) []item {
	var list []item
	ast.WalkFunc(node, func(node ast.Node, entering bool) ast.WalkStatus {
		switch v := node.(type) {
		case *ast.Text:
			list = append(list, item{
				Type:    "text",
				Content: string(v.Literal),
			})
		case *ast.CodeBlock:
			list = append(list, item{
				Type:    "code",
				Content: string(v.Literal),
			})
		}
		return 0
	})
	return list
}
