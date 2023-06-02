[![Go Report Card](https://goreportcard.com/badge/github.com/ystyle/sqlmanager)](https://goreportcard.com/report/github.com/ystyle/sqlmanager)
[![License](https://img.shields.io/badge/license-MulanPSL2-blue.svg)](https://github.com/ystyle/sqlmanager/blob/master/LICENSE)
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

#### create a markdown in
> markdown content in sql/test.md
```markdown
    ### GetStudentByID
    >get student by id, required id
    ```sql
    select * from student where id = {{.}}
    ```
```

#### Use Available plugins
1. File Driver: using *.md file
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
       sql, err := sm.RenderTPL("test/GetStudentByID", 1)
       if err != nil {
           panic(err)
       }
       fmt.Println(sql)
       // select * from student where id = 1
   
       sql, err = sm.RenderTPL("test2/GetStudentByID2", 1)
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
2. Embed markdown Driver: using go embed.
   ```go
   //go:embed test-sql
   var Assets embed.FS
   func main() {
       sm = sqlmanager.New()
       sm.Use(sqlmanager.NewMarkdownDriverWithEmbedDir(Assets, "test-sql"))
       // sm.Use(sqlmanager.NewMarkdownDriverWithEmbed(Assets)) // default dir is sql
       sm.load()
       sql, _ := sm.RenderTPL("test/GetStudentByID", 1)
   }
   ```
   >when the sql in `test-sql/admin/report.md/###GetStudentByRoot` the sql id is: `admin/report/GetStudentByRoot`

3.  Dynamic Driver： using variable.
    ```go
    store := sqlmanager.NewDynamicDriver()

    const GetStudentByID = "select * from student where id = {{.}}"
    store.RegisterWithDescs("GetStudentByID", "Query Student by ID", GetStudentByID)

    sm := sqlmanager.New()
    sm.Use(store)
    sm.load()
    sql, _ := sm.RenderTPL("GetStudentByID", 1)
    ```
4. Database Driver: store in database
   >you can custom the table name, and delete row should set deleted = 1, deleted_at = now(). you can build a interface to manage your sqls in your product.  
   >Field names: name, deleted, deleted_at, description, sql are required in this table.
   >load sql only on call `sm.load()`, it mean you should run `sm.load()` after changes
   ```sql
    create table sql_manager
    (
    id            int unsigned auto_increment primary key,
    name        varchar(255)  null,
    deleted     int default 0 not null,
    deleted_at  datetime      null,
    description varchar(255)  null,
    `sql`       text          null,
    constraint sql_manager_deleted_name_uindex
    unique (deleted, name)
    );
    INSERT INTO sql_manager (name, deleted, deleted_at, description, `sql`) VALUES ('GetStudentByID', 0, null, 'get student by id, required id', 'select * from student where id = {{.}}');
    ```
    ```go
    import (
       "database/sql"
        _ "github.com/go-sql-driver/mysql"
    )
    
    func main() {
        db, _ := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test1")
        sm = sqlmanager.New()
        sm.Use(sqlmanager.NewDatabaseDriver(db, "sql_manager"))
        sm.load()
        sql, _ := sm.RenderTPL("test/GetStudentByID", 1)
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

func (mdd *CustomeDriver ) DriverName() string {
return "CustomeDriver"
}

func (mdd *CustomeDriver ) Load() ([]sqlmanager.SqlTemple, error) {
var list []sqlmanager.SqlTemple
// db.table("sql_store").Find(&list)
return list, nil
}
```
