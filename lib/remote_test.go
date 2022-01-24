package githuborg

import (
	"os"
	"testing"

	"github.com/google/go-github/github"
)

func getOrg() GithubOrg {
	token := os.Getenv("GITHUB_TOKEN")
	return GetClient(token, "gitops-toolbox")
}

func addressOf(s string) *string {
	return &s
}

func TestGetRepository(t *testing.T) {
	var repo *github.Repository
	org := getOrg()
	expectedName := "test"
	repo, err := org.GetRepository("test")

	if err != nil {
		t.Error(err)
	}

	if *repo.Name != expectedName {
		t.Errorf("Name is %s, expected to be %s", *repo.Name, expectedName)
	}
}

func TestGetRepo(t *testing.T) {
	var repo *Repo
	org := getOrg()
	expectedName := "test"
	repo, err := org.GetRepo("test")

	if err != nil {
		t.Error(err)
	}

	if *repo.Name != expectedName {
		t.Errorf("Name is %s, expected to be %s", *repo.Name, expectedName)
	}
}
