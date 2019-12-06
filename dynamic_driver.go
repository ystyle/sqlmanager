package sqlmanager

type DynamicDriver struct {
	sqls []SqlTemple
}

func NewDynamicDriver() *DynamicDriver {
	return &DynamicDriver{}
}

func (dd *DynamicDriver) DriverName() string {
	return "Dynamic"
}

func (dd *DynamicDriver) Register(name, sql string) {
	dd.RegisterWithDescs(name, "", sql)
}

func (dd *DynamicDriver) RegisterWithDescs(name, description, sql string) {
	dd.sqls = append(dd.sqls, SqlTemple{
		Name:        name,
		Description: description,
		Sql:         sql,
	})
}

func (dd *DynamicDriver) Load() ([]SqlTemple, error) {
	return dd.sqls, nil
}
