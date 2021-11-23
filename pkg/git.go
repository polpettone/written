package pkg

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)


func GetHistory(repoPath, filePath string) (string, error) {
	repo, err := git.PlainOpen(repoPath)

	ref, err := repo.Head()
	cIter, err := repo.Log(&git.LogOptions{
		From: ref.Hash(),
		FileName: &filePath,

	})
	if err != nil {
		return "", err
	}
	history := ""

	err = cIter.ForEach(func(c *object.Commit) error {
		history += fmt.Sprintf("%s %s\n", c.Author.When, c.Message)
		return nil
	})

	return history, nil
}