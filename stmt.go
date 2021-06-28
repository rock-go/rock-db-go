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
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	st := &luaStmt{s:s , d:sqlDB}
	st.initMeta()
	L.Push(L.NewAnyData(st))
	return 1
}

func (s *luaStmt) initMeta() {
	s.meta.Set("Query" , lua.NewFunction(s.query))
	s.meta.Set("Exec" , lua.NewFunction(s.exec))
	s.meta.Set("Close" , lua.NewFunction(s.close))
}

func (s *luaStmt) query(L *lua.LState) int {
	args := getSTMTArgs(L)
	sqlRows, err := s.s.Query(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer sqlRows.Close()
	rows, columns, err := parseRows(sqlRows, L)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result := lua.NewUserKV()
	result.Set(`rows`, rows)
	result.Set(`columns`, columns)
	L.Push(result)
	return 1
}

func (s *luaStmt) exec(L *lua.LState) int {
	args := getSTMTArgs(L)
	sqlResult, err := s.s.Exec(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result := L.NewTable()
	if id, err := sqlResult.LastInsertId(); err == nil {
		result.RawSetString(`last_insert_id`, lua.LNumber(id))
	}
	if aff, err := sqlResult.RowsAffected(); err == nil {
		result.RawSetString(`rows_affected`, lua.LNumber(aff))
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

func getSTMTArgs(L *lua.LState) []interface{} {
	args := make([]interface{}, 0)
	for i := 2; i <= L.GetTop(); i++ {
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