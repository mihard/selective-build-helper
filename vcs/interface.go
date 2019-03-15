package vcs

type VCS interface {
	GetLastCommitID() (string, error)
	GetListOfChangedFiles(commit string) ([]string, error)
}
