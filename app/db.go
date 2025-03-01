package main

type ValExp struct {
	Val string
	Exp int64
}

var (
	setDB   = &(map[string]string{})
	setEXDB = &(map[string]ValExp{})
)

func SETDB_get() *map[string]string {
	return setDB
}

func SETEXDB_get() *map[string]ValExp {
	return setEXDB
}
