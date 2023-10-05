package backend

import (
	"github.com/charmbracelet/soft-serve/git"
	"github.com/charmbracelet/soft-serve/server/proto"
)

// LatestFile returns the contents of the latest file at the specified path in
// the repository and its file path.
func LatestFile(r proto.Repository, ref *git.Reference, pattern string) (string, string, error) {
	repo, err := r.Open()
	if err != nil {
		return "", "", err
	}
	return git.LatestFile(repo, ref, pattern)
}

// Readme returns the repository's README.
func Readme(r proto.Repository, ref *git.Reference) (readme string, path string, err error) {
	pattern := "[rR][eE][aA][dD][mM][eE]*"
	readme, path, err = LatestFile(r, ref, pattern)
	return
}
