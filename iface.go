// Package db implements golang package db functionality for lua.
package db

import (
	lua "github.com/rock-go/rock/lua"
)


type dbx interface {
	open(*config) (dbx , error)
	//Query(L *lua.LState)   int
	//Exec(L *lua.LState)    int
	//Stmt(L *lua.LState)    int
	//Command(L *lua.LState) int
	Stop(L *lua.LState)    int
}

func Open(L *lua.LState) int {

	cfg := newConfig(L)
	if e := cfg.verify() ; e != nil {
		L.RaiseError("%v" , e)
		return 0
	}

	create, ok := knownDrivers[cfg.Driver]
	if !ok {
		L.RaiseError("unkown Driver: %s" , cfg.Driver)
		return 0
	}

	db, err := create(cfg)
	if err != nil {
		L.RaiseError("%v" , err)
		return 0
	}

	L.Push(L.NewAnyData(db))
	return 1
}