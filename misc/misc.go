package misc

import (
	"strings"
)

func Imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func SetDebugFlags(flags string) {
	dbs := strings.Split(flags, ",")
	for _, db := range dbs {
		Db[db] = true
	}
}

var Db map[string]bool

func init() {
	Db = make(map[string]bool)
	Db["fx1"] = false // Output from processing of functions like __include__
}

/* vim: set noai ts=4 sw=4: */
