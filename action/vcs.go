package action

import (
	"fmt"
	"github.com/mihard/selective-build-helper/vcs"
	"github.com/pkg/errors"
	"log"
	"path/filepath"
	"regexp"
)

func CollectDirectories(bp string, commit string, vc vcs.VCS) (directories []string, err error) {
	if commit == "" {
		commit, err = vc.GetLastCommitID()
		if err != nil {
			return []string{}, errors.Wrapf(err, "Unable to fetch current commit message")
		}
	}

	log.Printf("Commit id: %s", commit)
	files, err := vc.GetListOfChangedFiles(commit)
	if err != nil {
		return []string{}, errors.Wrapf(err, "Unable to fetch list of changed files")
	}

	s := string(filepath.Separator)

	rxp := fmt.Sprintf("%s%s([^%s]+)%s", bp, s, s, s)
	rx, err := regexp.Compile(rxp)
	if err != nil {
		return []string{}, errors.Wrapf(err, "Unable to compile regexp")
	}

	uniqueFolders := map[string]string{}

	for _, f := range files {
		log.Printf("Changed file: %s", f)

		matches := rx.FindStringSubmatch(f)
		if len(matches) == 2 {
			uniqueFolders[matches[1]] = matches[1]
		}
	}

	for _, p := range uniqueFolders {
		directories = append(directories, p)
	}

	return directories, nil
}
