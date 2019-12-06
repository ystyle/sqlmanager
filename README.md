### sqlmanager
a library for manager sql with markdown like [beetsql](https://gitee.com/xiandafu/beetlsql)

**Not A Go ORM Library**
 
### feature
- manage sql with markdown
- render sql with go template
- support load sql with custom plugin

### usage

create a markdown in `sql/test.md`
```markdown
    ### GetStudentByID
    >get student by id, required id
    ```sql
    select * from student where id = {{.}}
    ```
```

in golang 
```go
package main

import (
    "fmt"
    "github.com/ystyle/sqlmanager"
)

func main() {
    sm := sqlmanager.New()
    sm.Use(NewMarkdownDriver())
    // sm.Use(NewMarkdownDriverWithDir("./prod-sql"))
    sm.Load()
    sql, err := sm.RenderTPL("GetStudentByID2", 1)
    if err != nil {
        panic(err)
    }
    fmt.Println(sql)
    // select * from student where id = 1
}
```

### custom puglin
> implement sqlmanager.Driver
```go
type CustomeDriver struct {
}

func NewCustomeDriver() *CustomeDriver {
	return &CustomeDriver{}
}

func (mdd *MarkdownDriver) DriverName() string {
	return "CustomeDriver"
}

func (mdd *MarkdownDriver) Load() ([]SqlTemple, error) {
    var list []SqlTemple
    // db.table("sql_store").Find(&list)
    return list, nil
}
```