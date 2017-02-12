package storage

import (
	"fmt"

	"remy.io/memoiz/log"
)

// InClause generates the IN clause such as,
// calling InClause(3,4) returns:
// ($3,$4,$5,$6)
func InClause(start, count int) string {
	if count == 0 {
		log.Error("InClause: called with c == 0")
		return ""
	}

	str := "("
	for i := 0; i < count; i++ {
		str = fmt.Sprintf("%s$%d", str, i+start)
		if i != count-1 {
			str += ","
		}
	}
	str += ")"

	return str
}

func Values(d ...interface{}) []interface{} {
	rv := make([]interface{}, 0)
	for _, v := range d {
		rv = append(rv, v)
	}
	return rv
}

func BasicLike(str string) string {
	if len(str) == 0 {
		return ""
	}
	return "%" + str + "%"
}

func FuzzyLike(str string) string {
	if len(str) == 0 {
		return ""
	}

	rv := ""
	for _, k := range str {
		rv += "%" + string(k)
	}
	return rv + "%"
}
