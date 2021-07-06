package db

import (
	"database/sql"
	"sync"

	_ "github.com/lib/pq"
	"github.com/rock-go/rock/lua"
)

type luaPG struct {
	cfg *config
	sync.Mutex
	db *sql.DB
	opts *sql.TxOptions
}

var (
	sharedPG     = make(map[string]*luaPG, 0)
	sharedPGLock = &sync.Mutex{}
)

func newLuaPG(cfg *config) (dbx, error) {
	if !cfg.Shared {
		goto OPEN
	}

	sharedPGLock.Lock()
	defer sharedPGLock.Unlock()

	if result, ok := sharedPG[cfg.Url]; ok {
		return result , nil
	}

OPEN:
	db := &luaPG{}
	return db.open(cfg)
}

func (pg *luaPG) open(cfg *config) (dbx , error) {
	db, err := sql.Open(`postgres`, cfg.Url)
	if err != nil {
		return nil, err
	}

	result := &luaPG{cfg: cfg}
	db.SetMaxIdleConns(cfg.Max)
	db.SetMaxOpenConns(cfg.Max)
	result.db = db

	if cfg.Shared {
		sharedPG[cfg.Url] = result
	}

	pg.opts = &sql.TxOptions{ReadOnly: cfg.ReadOnly}
	return result, nil

}

func (pg *luaPG) Stop(L *lua.LState) int {
	pg.Lock()
	defer pg.Unlock()
	err := pg.db.Close()
	if err != nil {
		L.Push(lua.LString(err.Error()))
	}
	if pg.cfg.Shared {
		sharedPGLock.Lock()
		delete(sharedPG, pg.cfg.Url)
		sharedPGLock.Unlock()
	}

	L.Push(lua.LNil)
	return 1
}


//func (pg *luaPG) Query(L *lua.LState)   int { return query(L , pg.db , pg.opts) }
//func (pg *luaPG) Exec(L *lua.LState)    int { return exec(L , pg.db , pg.opts)  }
//func (pg *luaPG) Command(L *lua.LState) int { return command(L , pg.db)            }
func (pg *luaPG) Stmt(L *lua.LState)    int { return newLuaStmt(L , pg.db)         }