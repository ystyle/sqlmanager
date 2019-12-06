package sqlmanager

import (
	"bytes"
	"log"
	"text/template"
)

func (sm *SqlManager) RenderTPL(name string, data interface{}) (string, error) {
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

func (sm *SqlManager) RenderTPLUnSave(name string, data interface{}) string {
	sql, err := sm.RenderTPL(name, data)
	if err != nil {
		log.Printf("sqlmanager - ERROR: %s, render error: %s", name, err.Error())
		return ""
	}
	return sql
}

func (sm *SqlManager) RenderTPLString(sql string, data interface{}) (string, error) {
	t := template.New("sql")
	_, err := t.Parse(sql)
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

func (sm *SqlManager) RenderTPLStringUnsave(sql string, data interface{}) string {
	sql, err := sm.RenderTPL(sql, data)
	if err != nil {
		log.Printf("sqlmanager - ERROR: render error: %s", err.Error())
		return ""
	}
	return sql
}
