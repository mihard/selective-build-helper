package vcs

type VCS interface {
	GetLastCommit() (*Commit, error)
	GetCommitData(commit string) (*Commit, error)
	GetListOfChangedFiles(commit *Commit) ([]string, error)
}

type Commit struct {
	ID        string
	ParentIDs []string
	Subject   string
}
