package db

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rock-go/rock/lua"
)

type luaMySQL struct {
	cfg *config
	sync.Mutex

	db *sql.DB
	opts *sql.TxOptions
}

var (
	sharedMySQL     = make(map[string]*luaMySQL, 0)
	sharedMySQLLock = &sync.Mutex{}
)

func (mysql *luaMySQL) open(cfg *config) (dbx , error) {
	db, err := sql.Open(`mysql`, cfg.Url)
	if err != nil {
		return nil, err
	}

	result := &luaMySQL{cfg: cfg}
	db.SetMaxIdleConns(cfg.Max)
	db.SetMaxOpenConns(cfg.Max)
	result.db = db

	if cfg.Shared {
		sharedMySQL[cfg.Url] = result
	}

	mysql.opts = &sql.TxOptions{ReadOnly: cfg.ReadOnly}
	return result, nil
}

func newLuaMysql(cfg *config) (dbx, error) {
	if !cfg.Shared {
		goto OPEN
	}

	sharedMySQLLock.Lock()
	defer sharedMySQLLock.Unlock()

	if result, ok := sharedMySQL[cfg.Url]; ok {
		return result , nil
	}

OPEN:
	db := &luaMySQL{}
	return db.open(cfg)
}

func (mysql *luaMySQL) Stop(L *lua.LState) int {
	mysql.Lock()
	defer mysql.Unlock()
	err := mysql.db.Close()
	if err != nil {
		L.Push(lua.LString(err.Error()))
	}
	if mysql.cfg.Shared {
		sharedMySQLLock.Lock()
		delete(sharedMySQL, mysql.cfg.Url)
		sharedMySQLLock.Unlock()
	}

	L.Push(lua.LNil)
	return 1
}


func (mysql *luaMySQL) Query(L *lua.LState)   int { return query(L , mysql.db , mysql.opts) }
func (mysql *luaMySQL) Exec(L *lua.LState)    int { return exec(L , mysql.db , mysql.opts)  }
func (mysql *luaMySQL) Command(L *lua.LState) int { return command(L , mysql.db)            }
func (mysql *luaMySQL) Stmt(L *lua.LState)    int { return newLuaStmt(L , mysql.db)         }