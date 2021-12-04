package githuborg

import (
	"context"
	"fmt"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

type Repo struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Homepage    *string `json:"homepage,omitempty"`

	Private     *bool   `json:"private,omitempty"`
	Visibility  *string `json:"visibility,omitempty"`
	HasIssues   *bool   `json:"has_issues,omitempty"`
	HasProjects *bool   `json:"has_projects,omitempty"`
	HasWiki     *bool   `json:"has_wiki,omitempty"`
	IsTemplate  *bool   `json:"is_template,omitempty"`

	Topics   []string `json:"topics,omitempty"`
	Archived *bool    `json:"archived,omitempty"`
	Disabled *bool    `json:"disabled,omitempty"`

	// Creating an organization repository. Required for non-owners.
	TeamID *int64 `json:"team_id,omitempty"`

	AutoInit            *bool   `json:"auto_init,omitempty"`
	GitignoreTemplate   *string `json:"gitignore_template,omitempty"`
	LicenseTemplate     *string `json:"license_template,omitempty"`
	AllowSquashMerge    *bool   `json:"allow_squash_merge,omitempty"`
	AllowMergeCommit    *bool   `json:"allow_merge_commit,omitempty"`
	AllowRebaseMerge    *bool   `json:"allow_rebase_merge,omitempty"`
	AllowAutoMerge      *bool   `json:"allow_auto_merge,omitempty"`
	DeleteBranchOnMerge *bool   `json:"delete_branch_on_merge,omitempty"`
}

func (r Repo) String() string {
	return github.Stringify(r)
}

func RepositoryToRepo(repo *github.Repository) *Repo {
	return &Repo{
		Name:                repo.Name,
		Description:         repo.Description,
		Homepage:            repo.Homepage,
		Private:             repo.Private,
		Visibility:          repo.Visibility,
		HasIssues:           repo.HasIssues,
		HasProjects:         repo.HasProjects,
		HasWiki:             repo.HasWiki,
		IsTemplate:          repo.IsTemplate,
		Topics:              repo.Topics,
		Archived:            repo.Archived,
		Disabled:            repo.Disabled,
		TeamID:              repo.TeamID,
		AutoInit:            repo.AutoInit,
		GitignoreTemplate:   repo.GitignoreTemplate,
		LicenseTemplate:     repo.LicenseTemplate,
		AllowSquashMerge:    repo.AllowSquashMerge,
		AllowMergeCommit:    repo.AllowMergeCommit,
		AllowRebaseMerge:    repo.AllowRebaseMerge,
		AllowAutoMerge:      repo.AllowAutoMerge,
		DeleteBranchOnMerge: repo.DeleteBranchOnMerge,
	}
}

func GetClient(accessToken string) (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "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"},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), ctx
}

func GetRepo(organization string, reponame string) (*github.Repository, error) {
	client, ctx := GetClient("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")

	repo, _, err := client.Repositories.Get(ctx, organization, reponame)

	return repo, err
}

func GetRepos(organization string) ([]*github.Repository, error) {

	client, ctx := GetClient("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")

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, organization, nil)

	return repos, err
}

func (ds Repo) setArchived(s *github.Repository) {

}

func (ds Repo) reconcile(s *github.Repository, localconfig []Repo) {
	// check if repo is in local config, warn if not and return

	// if repo is archived remotelly we can skip it, we are not touching archived repos
	if *s.Archived {
		fmt.Println("Repo is archived remotely, software will ignore archived repos")
		return
	}

	// if repo is archived locally we should clean up the remote repo and archive it

	// if repo is not archived we should try to sync the repo

	if s.Archived != ds.Archived {
		ds.setArchived(s)
	}
}
