package vcs

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

type Git struct {
	basePath string
}

func MakeGit(rootPath string, basePath string) *Git {
	err := os.Chdir(rootPath)
	if err != nil {
		log.Printf("Unable to chdir: %s", err.Error())
	}

	return &Git{basePath: basePath}
}

func (g *Git) GetLastCommitID() (string, error) {
	cmd := exec.Command("git", "log", "-1", `--pretty=format:%H`)
	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (g *Git) GetListOfChangedFiles(commit string) ([]string, error) {
	cmd := exec.Command("git", "diff-tree", "--no-commit-id", "--name-only", "-r", commit)
	output, err := cmd.Output()

	if err != nil {
		return []string{}, err
	}

	var files []string
	_files := strings.Split(string(output), "\n")

	for _, f := range _files {
		if f != "" {
			files = append(files, f)
		}
	}

	return files, nil
}
