package sqlmanager

import (
	"bytes"
	"errors"
	"fmt"
)

func (sm *SqlManager) RenderTPL(name string, data interface{}) (string, []any, error) {
	var buff bytes.Buffer
	args, err := sm.tpl.ExecuteTemplate(&buff, name, data)
	if err != nil {
		sql, has := sm.findTpl(name)
		if has {
			return "", nil, fmt.Errorf("sqlmanager - ERROR: %s[%s] %w", name, sql.Description, err)
		}
		return "", nil, fmt.Errorf("sqlmanager - ERROR: %s %w", name, errors.New(fmt.Sprintf("template: %s no found", name)))
	}
	return buff.String(), args, nil
}
