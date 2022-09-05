package sqlmanager

import (
	"github.com/gomarkdown/markdown/ast"
)

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
