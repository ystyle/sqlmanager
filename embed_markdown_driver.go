package sqlmanager

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"io/fs"
	"log"
	"path"
	"strings"
)

type EmbedMarkdownDriver struct {
	fs  embed.FS
	dir string
}

func NewMarkdownDriverWithEmbedDir(fs embed.FS, dir string) *EmbedMarkdownDriver {
	return &EmbedMarkdownDriver{
		fs:  fs,
		dir: dir,
	}
}

func NewMarkdownDriverWithEmbed(fs embed.FS) *EmbedMarkdownDriver {
	return NewMarkdownDriverWithEmbedDir(fs, "sql")
}
func (mdd *EmbedMarkdownDriver) DriverName() string {
	return "embed"
}

func (mdd *EmbedMarkdownDriver) Load() ([]SqlTemple, error) {
	var sqls []SqlTemple
	err := fs.WalkDir(mdd.fs, mdd.dir, func(subpath string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		ext := path.Ext(subpath)
		if ext == ".md" || ext == ".markdown" {
			s, err := mdd.parseMarkdown(subpath)
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

func (mdd *EmbedMarkdownDriver) parseMarkdown(filename string) ([]SqlTemple, error) {
	var sqls []SqlTemple
	buf, err := mdd.fs.ReadFile(filename)
	if err != nil {
		log.Printf("sqlmanager - ERROR: %s loading failed...\n", filename)
		return nil, err
	}
	if bytes.ContainsRune(buf, '\r') {
		buf = NormalizeNewlines(buf)
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
		var tpl SqlTemple
		if list[i].Type == "text" && list[i+1].Type == "code" {
			tpl.Name = mdd.getName(filename, list[i].Content)
			tpl.Sql = list[i+1].Content
			sqls = append(sqls, tpl)
			i += 2
		} else if list[i].Type == "text" && list[i+1].Type == "text" && list[i+2].Type == "code" {
			tpl.Name = mdd.getName(filename, list[i].Content)
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

func (mdd *EmbedMarkdownDriver) getName(filename, code string) string {
	ext := path.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	base = strings.TrimPrefix(base, mdd.dir)
	base = strings.TrimPrefix(base, "/")
	return path.Join(strings.TrimSuffix(base, ext), code)
}
