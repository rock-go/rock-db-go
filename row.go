package db

import (
	"database/sql"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xcall"
	"github.com/lib/pq"
)

type row map[string]interface{}

type xRows struct {
	err error
	lua.NoReflect
	data *sql.Rows
	cols *lua.LTable
	rows []row

	size int
	meta lua.UserKV
}

func newLuaRows(L *lua.LState , data *sql.Rows , e error) int {
	x := &xRows{
		err : e,
		data: data,
		size: 0,
		cols: nil,
		meta: lua.NewUserKV(),
	}
	if e == nil {
		x.parse(L)
		defer data.Close()
	}

	L.Push(L.NewAnyData(x))
	return 1
}

func (r *row) DisableReflect() {}

func (r *row) Get(L *lua.LState , key string ) lua.LValue {
	val , ok := (*r)[key]
	if !ok {
		return lua.LNil
	}
	converted := val.([]uint8)

	arr := make([]string , 0)
	pqArr := pq.Array(&arr)
	err := pqArr.Scan(converted)
	if err != nil {
		return lua.LString(converted)
	}
	tab := L.NewTable()
	for _ , v := range arr {
		tab.Append(lua.LString(v))
	}

	return tab
}

func (r *row) Set(L *lua.LState , key string , val lua.LValue) {
	(*r)[key] = val
}

func (x *xRows) parse(L *lua.LState) {
	cols , err := x.data.Columns()
	if err != nil {
		x.err = err
		return
	}

	//同步列字段
	n := len(cols)
	tab := L.CreateTable(n , 0)
	for i := 0 ;i < n ; i++ {
		tab.Append(lua.LString(cols[i]))
	}
	x.cols = tab

	//解析列数据
	rows := make([]row , 0)
	items := make([]interface{} ,  n)
	ptrs := make([]interface{} , n)
	for i := range items {
		ptrs[i] = &items[i]
	}

	for x.data.Next() {

		x.size++
		err = x.data.Scan(ptrs...)
		if err != nil {
			x.err = err
			return
		}

		r := make(row , n)
		for i , key := range cols {
			 r[key] =  *(ptrs[i].(*interface{}))
		}
		rows = append(rows, r)
	}
	x.rows = rows
}

func (x *xRows) Get(L *lua.LState , key string ) lua.LValue {
	var lv lua.LValue
	if lv = x.meta.Get(key);lv != lua.LNil {
		return lv
	}

	switch key {

	//固定值
	case "cols":
		return x.cols

	case "ERR":
		if x.err == nil {
			return lua.LNil
		}
		return lua.LString(x.err.Error())

	case "size":
		return lua.LNumber(x.size)

	case "pairs":
		lv = lua.NewFunction(x.Pairs)

	case "try_catch":
		lv = lua.NewFunction(x.tryCatch)

	default:
		return lua.LNil
	}

	x.meta.Set(key , lv)
	return lv
}

func (x *xRows) Pairs(L *lua.LState) int {
	fn := L.CheckFunction(1)
	p := lua.P{
		Fn:fn,
		NRet: 0,
		Protect: true,
	}

	stop := false
	stopFn := newStopFn(&stop)

	for i := 0 ; i < x.size ; i++ {
		r := x.rows[i]
		err := xcall.CallByParam(L , p , xcall.Rock ,
			lua.LNumber(i + 1) , L.NewAnyData(&r) , stopFn)
		if err != nil {
			L.RaiseError("%v" , err)
			return  0
		}

		if stop {
			return 0
		}
	}

	return 0
}

func (x *xRows) tryCatch(L *lua.LState) int {
	if x.err == nil {
		return 0
	}
	L.RaiseError("%v" , x.err)
	return 0
}