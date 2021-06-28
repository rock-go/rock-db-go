package db

import (
	"context"
	"github.com/rock-go/rock/lua"
	"database/sql"
)

func query(L *lua.LState , sqlDB *sql.DB , opts *sql.TxOptions)  int {
	str := L.CheckString(1)
	tx, err := sqlDB.BeginTx(context.Background(), opts)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer tx.Commit()
	sqlRows, err := tx.Query(str)
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

	ret := lua.NewUserKV()
	ret.Set(`rows`, rows)
	ret.Set(`columns`, columns)
	L.Push(ret)
	return 1
}

// Exec lua db.exec(query) returns ({rows_affected=number, last_insert_id=number}, err)
func exec(L *lua.LState , sqlDB *sql.DB , opts *sql.TxOptions) int {
	str := L.CheckString(1)
	tx, err := sqlDB.BeginTx(context.Background(), opts)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer tx.Commit()
	sqlResult, err := tx.Exec(str)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result := lua.NewUserKV()
	if id, err := sqlResult.LastInsertId(); err == nil {
		result.Set(`last_insert_id`, lua.LNumber(id))
	}
	if aff, err := sqlResult.RowsAffected(); err == nil {
		result.Set(`rows_affected`, lua.LNumber(aff))
	}
	L.Push(result)
	return 1
}

// Command lua db_ud:command(query) returns ({rows = {}, columns = {}}, err)
func command(L *lua.LState , sqlDB *sql.DB)  int {
	str := L.CheckString(1)
	sqlRows, err := sqlDB.Query(str)
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