[![Go Report Card](https://goreportcard.com/badge/github.com/ystyle/sqlmanager)](https://goreportcard.com/report/github.com/ystyle/sqlmanager)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/ystyle/sqlmanager/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/ystyle/sqlmanager?status.svg)](https://godoc.org/github.com/ystyle/sqlmanager)
---
### sqlmanager
a library for manager sql with markdown or constant. and you can custom sql store plugin.

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
    sm.Use(sqlmanager.NewMarkdownDriver())
    // load sql with custom dir
    // sm.Use(sqlmanager.NewMarkdownDriverWithDir("./prod-sql"))
    // register go template func
    // sm.RegisterFunc(template.FuncMap{
    //     "test": func(v string) string {
    //         return strings.ToUpper(v)
    //     },
    // })
    sm.Load()
    sql, err := sm.RenderTPL("GetStudentByID", 1)
    if err != nil {
        panic(err)
    }
    fmt.Println(sql)
    // select * from student where id = 1
    
    // using gorm 
    // db, err := gorm.Open("databaseurl")
    // if err != nil {
    //   panic("failed to connect database")
    // }
    // db.Raw(sql)
    
    // using database/sql
    // db, err := sql.Open("driver-name", "database=test1")
    // if err != nil {
    //   log.Fatal(err)
    // }
    // db.Query(sql)
}
```

### Available plugins
- File Driver: using file, See above.
- Dynamic Driver： using variable.
```go
store := sqlmanager.NewDynamicDriver()

const GetStudentByID = "select * from student where id = {{.}}"
store.RegisterWithDescs("GetStudentByID", "Query Student by ID", GetStudentByID)

sm := sqlmanager.New()
sm.Use(store)
sm.load()
```
- Embed markdown Driver: using go embed.
```go
//go:embed test-sql
var Assets embed.FS
func main() {
	sm = sqlmanager.New()
    sm.Use(sqlmanager.NewMarkdownDriverWithEmbedDir(Assets. "test-sql"))
    // sm.Use(sqlmanager.NewMarkdownDriverWithEmbed(Assets)) // default dir is sql
    sm.load()
}
```
> when the sql in `test-sql/admin/report.md/###GetStudentByRoot` the sql id is: `admin.report.GetStudentByRoot`


### custom puglin
> implement sqlmanager.Driver
```go
type CustomeDriver struct {
}

func NewCustomeDriver() *CustomeDriver {
    return &CustomeDriver{}
}

func (mdd *CustomeDriver ) DriverName() string {
    return "CustomeDriver"
}

func (mdd *CustomeDriver ) Load() ([]sqlmanager.SqlTemple, error) {
    var list []sqlmanager.SqlTemple
    // db.table("sql_store").Find(&list)
    return list, nil
}
```
