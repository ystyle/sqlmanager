package sqlmanager

import (
	"bytes"
	"log"
	"text/template"
)

func (sm *sqlManager) RenderTPL(name string, data interface{}) (string, error) {
	tpl, err := sm.findTpl(name)
	if err != nil {
		return "", err
	}
	t := template.New(name)
	_, err = t.Parse(tpl.Sql)
	if err != nil {
		return "", err
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, data)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

func (sm *sqlManager) RenderTPLUnSave(name string, data interface{}) string {
	sql, err := sm.RenderTPL(name, data)
	if err != nil {
		log.Printf("sqlmanager - ERROR: %s, render error: %s", name, err.Error())
		return ""
	}
	return sql
}
