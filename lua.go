package db

import (
	lua "github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xcall"
)

func LuaInjectApi(env xcall.Env) {
	uv := lua.NewUserKV()
	uv.Set("open" , lua.NewFunction(Open))
	env.SetGlobal("db" , uv)
}
