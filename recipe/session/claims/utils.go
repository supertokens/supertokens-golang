package claims

import (
	"fmt"
)

func assertCondition(condition bool, message ...interface{}) {
	if !condition {
		panic(fmt.Sprint(message...))
	}
}

func includes(s []interface{}, e interface{}) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
