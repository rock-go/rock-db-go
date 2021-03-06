package db

import (
	"database/sql"

	lua "github.com/rock-go/rock/lua"
)

type luaStmt struct {
	lua.NoReflect

	s *sql.Stmt
	d *sql.DB
	meta lua.UserKV
}

func newLuaStmt(L *lua.LState , sqlDB *sql.DB) int {
	str := L.CheckString(1)

	s , err := sqlDB.Prepare(str)
	if err != nil {
		L.RaiseError("db stmt %v" , err)
		return 0
	}

	st := &luaStmt{s:s , d:sqlDB , meta: lua.NewUserKV()}
	L.Push(L.NewAnyData(st))
	return 1
}

func (s *luaStmt) query(L *lua.LState) int {
	args := getSTMTArgs(L)
	sqlRows, err := s.s.Query(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return newLuaRows(L , sqlRows , err)
}

func (s *luaStmt) exec(L *lua.LState) int {
	args := getSTMTArgs(L)
	sqlResult, err := s.s.Exec(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result := lua.NewUserKV()
	if id, e := sqlResult.LastInsertId(); e == nil {
		result.Set(`id`, lua.LNumber(id))
	}
	if aff, e := sqlResult.RowsAffected(); e == nil {
		result.Set(`affected`, lua.LNumber(aff))
	}
	L.Push(result)
	return 1
}

func (s *luaStmt) close(L *lua.LState) int {
	if err := s.s.Close(); err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

func (s *luaStmt) Get(L *lua.LState , key string) lua.LValue {

	if lv := s.meta.Get(key); lv != lua.LNil {
		return lv
	}

	var fn *lua.LFunction
	switch key {
	case "Exec":
		fn = L.NewFunction(s.exec)
	case "Query":
		fn = L.NewFunction(s.query)
	case "Close":
		fn = L.NewFunction(s.close)
	default:
		return lua.LNil
	}

	s.meta.Set(key , fn)
	return fn
}

func getSTMTArgs(L *lua.LState) []interface{} {
	args := make([]interface{}, 0)
	for i := 1; i <= L.GetTop(); i++ {
		any := L.CheckAny(i)
		switch any.Type() {
		case lua.LTNil:
			args = append(args, nil)
		default:
			args = append(args, L.CheckAny(i))
		}
	}
	return args
}