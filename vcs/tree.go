package vcs

import "github.com/mihard/selective-build-helper/func/slices"

type CommitTree [][]string

func (t CommitTree) AsPlainSlice() (plain []string) {
	for _, layer := range t {
		plain = append(plain, layer...)
	}

	return slices.UniqueStrings(plain)
}

func (t CommitTree) GetDeepestLayer() (deepest []string) {
	return t[len(t)-1]
}
