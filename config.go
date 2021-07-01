package db

import (
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xreflect"
)

type config struct {
	Driver   string `lua:"driver"     type:"string"`
	Url      string `lua:"url"        type:"string"`
	Shared   bool   `lua:"share,true" type:"bool"`
	Max      int    `lua:"Max,1"      type:"int"`
	ReadOnly bool   `lua:"readonly,false" type:"bool"'`
}

func newConfig(L *lua.LState) *config {
	tab := L.CheckTable(1)
	cfg := &config{}
	if e := xreflect.ToStruct(tab , cfg); e != nil {
		L.RaiseError("%v" , e)
		return nil
	}
	return cfg
}

func (cfg *config) verify() error {
	return nil
}