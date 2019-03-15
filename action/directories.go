package action

import (
	"io/ioutil"
	"path/filepath"
)

func CollectAllDirectories(rp string, bp string) (list []string, err error) {
	files, err := ioutil.ReadDir(filepath.Join(rp, bp))
	if err != nil {
		return []string{}, err
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		list = append(list, f.Name())
	}

	return list, nil
}
