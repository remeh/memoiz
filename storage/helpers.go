package storage

import (
	"fmt"

	"remy.io/memoiz/log"
)

// InClause generates the IN clause such as,
// calling with 2 returns:
// ($1,$2)
func InClause(c int) string {
	if c == 0 {
		log.Error("InClause: called with c == 0")
		return ""
	}

	str := "("
	for i := 0; i < c; i++ {
		str = fmt.Sprintf("%s ?%d", str, i+1)
		if i != c-1 {
			str += ","
		}
	}
	str += ")"

	return str
}

func Values(d ...interface{}) []interface{} {
	rv := make([]interface{}, len(d))
	for i, v := range d {
		rv[i] = v
	}
	return rv
}
