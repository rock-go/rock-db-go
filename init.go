package db

const (
	// Max open connections
	MaxOpenConns = 1
)

var (
	knownDrivers     = make(map[string]createFn, 2)
)

type createFn func(*config) (dbx , error)

func init() {
	knownDrivers["mysql"] = newLuaMysql
	knownDrivers["postgres"] = newLuaPG
}