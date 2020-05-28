package slices

import "sort"

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

func UniqueStrings(array []string) (r []string) {
	um := make(map[string]int)
	r = []string{}
	var rr []int

	for i, v := range array {
		um[v] = i
	}

	for _, i := range um {
		rr = append(rr, i)
	}

	sort.Ints(rr)

	for _, v := range rr {
		r = append(r, array[v])
	}

	return
}
