package githuborg

import (
	"context"
	"path/filepath"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func RepositoryToRepo(repo *github.Repository) *Repo {
	return &Repo{
		Name:              repo.Name,
		Description:       repo.Description,
		Homepage:          repo.Homepage,
		Private:           repo.Private,
		HasIssues:         repo.HasIssues,
		HasProjects:       repo.HasProjects,
		HasWiki:           repo.HasWiki,
		Topics:            repo.Topics,
		Archived:          repo.Archived,
		TeamID:            repo.TeamID,
		AutoInit:          repo.AutoInit,
		GitignoreTemplate: repo.GitignoreTemplate,
		LicenseTemplate:   repo.LicenseTemplate,
		AllowSquashMerge:  repo.AllowSquashMerge,
		AllowMergeCommit:  repo.AllowMergeCommit,
		AllowRebaseMerge:  repo.AllowRebaseMerge,
	}
}

type GithubOrg struct {
	client       *github.Client
	ctx          context.Context
	organization string
}

func GetClient(accessToken string, organization string) GithubOrg {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return GithubOrg{
		client:       github.NewClient(tc),
		ctx:          ctx,
		organization: organization,
	}
}

func (org GithubOrg) SyncTopics(repo *github.Repository) error {
	_, _, err := org.client.Repositories.ReplaceAllTopics(org.ctx, org.organization, *repo.Name, repo.Topics)

	return err
}

func (org GithubOrg) Sync(repo *github.Repository) error {
	_, _, err := org.client.Repositories.Edit(org.ctx, org.organization, *repo.Name, repo)

	return err
}

func (org GithubOrg) GetOutOfSyncRepos() ([]Repo, error) {
	repos, err := LoadConfig(filepath.Join("repos", org.organization))

	result := []Repo{}
	if err != nil {
		return result, err
	}

	for _, r := range repos {
		synced, _, err := r.InSync(org)
		if err == nil && !synced {
			result = append(result, r)
		}
	}

	return result, nil
}
