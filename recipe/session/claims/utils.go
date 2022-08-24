package claims

func includes(s []interface{}, e interface{}) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func includesAll(s []interface{}, e []interface{}) bool {
	valsMap := map[interface{}]bool{}
	for _, v := range s {
		valsMap[v] = true
	}
	for _, v := range e {
		if !valsMap[v] {
			return false
		}
	}
	return true
}

func excludesAll(s []interface{}, e []interface{}) bool {
	valsMap := map[interface{}]bool{}
	for _, v := range s {
		valsMap[v] = true
	}
	for _, v := range e {
		if valsMap[v] {
			return false
		}
	}
	return true
}
