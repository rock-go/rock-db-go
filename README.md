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

## 模板 *.Stmt(string)
模板查询
```lua
    local s ,err = mysql.Stmt("select *from admin where name= ?")
    --err
    
    local r , err = s.Exec("admin")
    --err
```

## 查询 rows = stmt.Query(string)
直接查询
```lua
    local stmt = mysql.Stmt("select * from admin where 1=?")
    local rows = mysql.Query('1')
    rows.try_catch() --判断异常
    
    --遍历数据
    rows.pairs(function(i , row , stop) -- i:索引,row:单条数据 stop:函数
        if i == 1 then
            print(row.field)    
        end
        
        if row.name == "admin" then
            return stop()    
        end
        
    end)
```


## 执行 stmt.Exec(string)
其他语句执行
```lua
    local stmt = mysql.Stmt("insert into t values(?)")

    local _ , err = mysql.Exec("create table t (id int)")

    local r , err = mysql.Exec("insert into t values(1)")
    print(r.affected)
    print(r.id)
```