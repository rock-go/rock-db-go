package db

import (
	"github.com/rock-go/rock/lua"
	"bytes"
	"fmt"
	"github.com/lib/pq"
)

//func query(L *lua.LState , sqlDB *sql.DB , opts *sql.TxOptions)  int {
//	str := L.CheckString(1)
//	stmt , err := sqlDB.Prepare(str)
//	if err != nil {
//		return newLuaRows(L , nil , err)
//	}
//	defer stmt.Close()
//
//	sqlRows, err := stmt.Query(str)
//	if err != nil {
//		return newLuaRows(L , nil , err)
//	}
//	return newLuaRows(L , sqlRows , nil)
//}
//
//// Exec lua db.exec(query) returns ({rows_affected=number, last_insert_id=number}, err)
//func exec(L *lua.LState , sqlDB *sql.DB , opts *sql.TxOptions) int {
//	str := L.CheckString(1)
//	tx, err := sqlDB.BeginTx(context.Background(), opts)
//	if err != nil {
//		L.Push(lua.LNil)
//		L.Push(lua.LString(err.Error()))
//		return 2
//	}
//	defer tx.Commit()
//	sqlResult, err := tx.Exec(str)
//	if err != nil {
//		L.Push(lua.LNil)
//		L.Push(lua.LString(err.Error()))
//		return 2
//	}
//	result := lua.NewUserKV()
//	if id, err := sqlResult.LastInsertId(); err == nil {
//		result.Set(`id`, lua.LNumber(id))
//	}
//	if aff, err := sqlResult.RowsAffected(); err == nil {
//		result.Set(`affected`, lua.LNumber(aff))
//	}
//	L.Push(result)
//	return 1
//}
//
//// Command lua db_ud:command(query) returns ({rows = {}, columns = {}}, err)
//func command(L *lua.LState , sqlDB *sql.DB)  int {
//	str := L.CheckString(1)
//	sqlRows, err := sqlDB.Query(str)
//	if err != nil {
//		return newLuaRows(L , nil , err)
//	}
//	defer sqlRows.Close()
//	return newLuaRows(L , sqlRows , nil)
//}

func newStopFn(flag *bool) *lua.LFunction {
	return lua.NewFunction(func( _ *lua.LState) int {
		*flag = true
		return 0
	})
}

func WriteToJson( buff *bytes.Buffer, v []byte ) {
	n := len(v)
	var ch byte
	for i := 0;i< n;i++ {
		ch = v[i]
		switch ch {
		case '"':
			buff.WriteByte('\\')
			buff.WriteByte('"')
		case '\\':
			buff.WriteByte('\\')
			buff.WriteByte('\\')
		case '\r':
			buff.WriteByte('\\')
			buff.WriteByte('r')
		case '\n':
			buff.WriteByte('\\')
			buff.WriteByte('n')
		case '\t':
			buff.WriteByte('\\')
			buff.WriteByte('t')
		default:
			buff.WriteByte(ch)
		}
	}
}

func WriteKeyToJson(buff *bytes.Buffer , key string) {
	buff.WriteByte('"')
	WriteToJson(buff , lua.S2B(key))
	buff.WriteByte('"')
	buff.WriteByte(':')
}

func WriteStrToJson(buff *bytes.Buffer , v string) {
	buff.WriteByte('"')
	WriteToJson(buff , lua.S2B(v))
	buff.WriteByte('"')
}

func WriteIntToJson(buff *bytes.Buffer , v interface{}) {
	val := fmt.Sprintf("%d" , v)
	WriteToJson(buff , lua.S2B(val))
}
func WriteFloatToJson(buff *bytes.Buffer , v interface{}) {
	val := fmt.Sprintf("%f" , v)
	WriteToJson(buff , lua.S2B(val))
}

func toArr(v []uint8) ([]string , error) {
	arr := make([]string , 0)
	pqArr := pq.Array(&arr)
	err := pqArr.Scan(v)
	if err != nil {
		return nil , err
	}
	return arr , nil
}
