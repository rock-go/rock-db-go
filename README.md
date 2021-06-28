# rock-db-go
磐石框架数据库驱动

## db.open{cfg}
创建配置
```lua
    local mysql = db.open{
        driver = "mysql",
        url    = "root:pass123!@tcp(127.0.0.1:3306)/rock?charset=utf8&parseTime=True",
        max    = 1024,
        shared = true,
        readonly = true,
    }
```

## 查询 *.Query(string)
直接查询
```lua
    local r , err = mysql.Query("select * from admin")
    for k , v in pairs(r.row) do
        --todo
    end
```

## 模板 *.Stmt(string)
模板查询
```lua
    local s ,err = mysql.Stmt("select *from admin where name= ?")
    --err
    
    local r , err = s.Exec("admin")
    --err
```

## 执行 *.Exec(string)
其他语句执行
```lua
    local _ , err = mysql.Exec("create table t (id int)")

    local r , err = mysql.Exec("insert into t values(1)")
    print(r.rows_affected)
```
