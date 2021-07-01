package db

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rock-go/rock/lua"
)

type luaMySQL struct {
	sync.RWMutex
	cfg *config
	db *sql.DB
	opts *sql.TxOptions
}

var (
	sharedMySQL     = make(map[string]*luaMySQL, 0)
	sharedMySQLLock = &sync.RWMutex{}
)

func newLuaMysql(cfg *config) (dbx, error) {
	var db dbx
	var ok bool

	if !cfg.Shared {
		goto OPEN
	}

	sharedMySQLLock.Lock()
	db , ok = sharedMySQL[cfg.Url]
	sharedMySQLLock.Unlock()
	if ok {
		return db , nil
	}

OPEN:
	db = &luaMySQL{}
	return db.open(cfg)
}

func getLuaMysql(L *lua.LState) int {
	url := L.CheckString(1)

	sharedMySQLLock.RLock()
	db , ok := sharedMySQL[url]
	sharedMySQLLock.RUnlock()
	if ok {
		L.Push(L.NewAnyData(db))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

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
		sharedMySQLLock.Lock()
		sharedMySQL[cfg.Url] = result
		sharedMySQLLock.Unlock()
	}

	mysql.opts = &sql.TxOptions{ReadOnly: cfg.ReadOnly}
	return result, nil
}

func (mysql *luaMySQL) Close() error {
	if mysql.cfg.Shared {
		sharedMySQLLock.Lock()
		delete(sharedMySQL, mysql.cfg.Url)
		sharedMySQLLock.Unlock()
	}

	mysql.Lock()
	err := mysql.db.Close()
	mysql.Unlock()
	return err

}

func (mysql *luaMySQL) Stop(L *lua.LState) int {
	err := mysql.Close()
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

func (mysql *luaMySQL) Query(L *lua.LState)   int {
	mysql.RLock()
	ret := query(L , mysql.db , mysql.opts)
	mysql.RUnlock()
	return ret
}
func (mysql *luaMySQL) Exec(L *lua.LState)    int {
	mysql.RLock()
	ret := exec(L , mysql.db , mysql.opts)
	mysql.RUnlock()
	return ret
}
func (mysql *luaMySQL) Command(L *lua.LState) int {
	mysql.RLock()
	ret := command(L , mysql.db)
	mysql.RUnlock()
	return ret
}
func (mysql *luaMySQL) Stmt(L *lua.LState)    int {
	mysql.RLock()
	ret := newLuaStmt(L , mysql.db)
	mysql.RUnlock()
	return ret
}