package slices

func InStrings(val string, array []string) (exists bool) {
	exists = false

	for _, v := range array {
		if val == v {
			exists = true
			return
		}
	}

	return
}

func UniqueStrings(input []string) (unique []string) {
	set := make(map[string]struct{})
	for _, item := range input {
		if _, found := set[item]; !found {
			set[item] = struct{}{}
			unique = append(unique, item)
		}
	}
	return
}
