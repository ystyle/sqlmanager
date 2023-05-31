package sqlmanager

import (
	"embed"
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
	buf, err := mdd.fs.ReadFile(filename)
	if err != nil {
		log.Printf("sqlmanager - ERROR: %s loading failed...\n", filename)
		return nil, err
	}
	return parseMarkdown(buf, mdd.getName(filename))
}

func (mdd *EmbedMarkdownDriver) getName(filename string) string {
	ext := path.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	base = strings.TrimPrefix(base, mdd.dir)
	base = strings.TrimPrefix(base, "/")
	return base
}
