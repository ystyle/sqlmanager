package sqlmanager

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
)

func (sm *SqlManager) RenderTPL(name string, data interface{}) (string, error) {
	tpl, err := sm.findTpl(name)
	if err != nil {
		return "", fmt.Errorf("sqlmanager - ERROR: %s %w", name, err)
	}
	t := template.New(name)
	_, err = t.Parse(tpl.Sql)
	if err != nil {
		return "", fmt.Errorf("sqlmanager - ERROR: %s[%s] %w", name, tpl.Description, err)
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, data)
	if err != nil {
		return "", fmt.Errorf("sqlmanager - ERROR: %s[%s] %w", name, tpl.Description, err)
	}
	return buff.String(), nil
}

func (sm *SqlManager) RenderTPLUnSave(name string, data interface{}) string {
	sql, err := sm.RenderTPL(name, data)
	if err != nil {
		log.Printf("sqlmanager - ERROR: %s, render error: %s", name, err.Error())
		return ""
	}
	return sql
}
