package vcs

import (
	"errors"
	"github.com/mihard/selective-build-helper/func/slices"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Git struct {
	basePath string
	logRX    *regexp.Regexp
}

const (
	head     = "HEAD"
	maxDepth = 100
)

func MakeGit(rootPath string, basePath string) *Git {
	err := os.Chdir(rootPath)
	if err != nil {
		log.Printf("Unable to chdir: %s", err.Error())
	}

	return &Git{
		basePath: basePath,
		logRX:    regexp.MustCompile("^((?:\\s*\\w{40}){2,})\\s+(.+)$"),
	}
}

func (g *Git) GetCommitData(commit string) (*Commit, error) {
	cmd := exec.Command("git", "show", commit, "-s", `--pretty=format:%H %P %s`)
	output, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	parts := g.logRX.FindStringSubmatch(string(output))

	if len(parts) != 3 {
		return nil, errors.New("unable to read git show")
	}

	ids := strings.Split(parts[1], " ")

	if len(ids) < 2 {
		return nil, errors.New("unable to read git show")
	}

	return &Commit{
		ID:        ids[0],
		ParentIDs: ids[1:],
		Subject:   parts[2],
	}, nil
}

func (g *Git) GetLastCommit() (*Commit, error) {
	return g.GetCommitData(head)
}

func (g *Git) GetListOfChangedFiles(c *Commit) ([]string, error) {
	if len(c.ParentIDs) > 1 {
		return g.collectMergedCommits(c)
	}

	return g.collectChangedFilesByCommitID(c.ID)
}

func (g *Git) collectMergedCommits(c *Commit) ([]string, error) {
	main := c.ParentIDs[0]
	merged := c.ParentIDs[1]

	commits := []string{c.ID}
	commits = append(commits, g.findAllMergedCommits(main, merged)...)

	var files []string

	for _, commit := range commits {
		fls, err := g.collectChangedFilesByCommitID(commit)
		if err != nil {
			log.Printf("WARN unable to collect files from merged PR: %s", err.Error())
		}

		for _, f := range fls {
			files = append(files, f)
		}
	}

	return slices.UniqueStrings(files), nil
}

func (g *Git) findAllMergedCommits(main string, merged string) []string {
	mainLine := CommitTree{[]string{main}}
	mainLine = g.loadMoreMailBranchCommits(mainLine)

	allMerged := []string{}

	toCheck := CommitTree{[]string{merged}}

	depth := maxDepth
	for depth >= 0 {
		plainMailLine := mainLine.AsPlainSlice()

		var unmerged []string

		for _, commit := range toCheck.GetDeepestLayer() {
			if !slices.InStrings(commit, plainMailLine) {
				unmerged = append(unmerged, commit)
			}
		}

		if len(unmerged) < 1 {
			return allMerged
		}

		allMerged = append(allMerged, unmerged...)

		var nextLayer []string

		for _, commit := range toCheck.GetDeepestLayer() {
			cd, err := g.GetCommitData(commit)
			if err != nil {
				log.Printf("ERROR unable to collect merged commit data: %s", err.Error())
				return allMerged
			}

			nextLayer = append(nextLayer, cd.ParentIDs...)
		}

		toCheck = append(toCheck, nextLayer)

		mainLine = g.loadMoreMailBranchCommits(mainLine)

		depth--
	}

	return allMerged
}

func (g *Git) collectChangedFilesByCommitID(id string) ([]string, error) {
	cmd := exec.Command("git", "diff-tree", "--no-commit-id", "--name-only", "-r", id)
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

func (g *Git) loadMoreMailBranchCommits(tree CommitTree) CommitTree {
	var nextLayer []string

	for _, commit := range tree.GetDeepestLayer() {
		cd, err := g.GetCommitData(commit)

		if err != nil {
			log.Printf("ERROR unable to collect merged commit data: %s", err.Error())
			return tree
		}

		for _, p := range cd.ParentIDs {
			nextLayer = append(nextLayer, p)
		}
	}

	return append(tree, nextLayer)
}
