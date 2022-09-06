package sqlmanager

import (
	"bytes"
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

func NormalizeNewlines(d []byte) []byte {
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = bytes.Replace(d, []byte{13, 10}, []byte{10}, -1)
	// replace CF \r (mac) with LF \n (unix)
	d = bytes.Replace(d, []byte{13}, []byte{10}, -1)
	return d
}
