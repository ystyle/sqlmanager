package sqlmanager

import (
	"bytes"
	"errors"
	"fmt"
	"log"
)

func (sm *SqlManager) RenderTPL(name string, data interface{}) (string, error) {
	var buff bytes.Buffer
	err := sm.tpl.ExecuteTemplate(&buff, name, data)
	if err != nil {
		sql, has := sm.findTpl(name)
		if has {
			return "", fmt.Errorf("sqlmanager - ERROR: %s[%s] %w", name, sql.Description, err)
		}
		return "", fmt.Errorf("sqlmanager - ERROR: %s %w", name, errors.New(fmt.Sprintf("template: %s no found", name)))
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
