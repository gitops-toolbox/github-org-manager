package githuborg

import (
	"fmt"
	"reflect"
	"sort"

	"log"

	"github.com/google/go-github/github"
)

type Team struct {
	Name       *string `json:"name, omitempty"`
	Slug       *string `json:"slug,omitempty"`
	ID         *int64  `json:"id,omitempty"`
	Permission *string `json:"permission,omitempty"`
}

type Collaborator struct {
	Login *string `json:"login,omitempty"`
	Name  *string `json:"name,omitempty"`

	// Permissions identifies the permissions that a user has on a given
	// repository. This is only populated when calling Repositories.ListCollaborators.
	Permissions map[string]bool `json:"permissions,omitempty"`
}

type Repo struct {
	Name          *string `json:"name,omitempty"`
	Description   *string `json:"description,omitempty"`
	Homepage      *string `json:"homepage,omitempty"`
	DefaultBranch *string `json:"default_branch,omitempty"`

	Private     *bool   `json:"private,omitempty"`
	Visibility  *string `json:"visibility,omitempty"`
	HasIssues   *bool   `json:"has_issues,omitempty"`
	HasProjects *bool   `json:"has_projects,omitempty"`
	HasWiki     *bool   `json:"has_wiki,omitempty"`
	IsTemplate  *bool   `json:"is_template,omitempty"`

	Teams         []*Team         `json:"teams,omitempty"`
	Collaborators []*Collaborator `json:"collaborators,omitempty"`
	Topics        []string        `json:"topics,omitempty"`
	Archived      *bool           `json:"archived,omitempty"`
	Disabled      *bool           `json:"disabled,omitempty"`

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

func (r Repo) InSync(org GithubOrg) (bool, *github.Repository, error) {
	repo, err := org.GetRepository(*r.Name)

	if err != nil {
		fmt.Println(err)
		return false, nil, err
	}

	sort.Strings(r.Topics)
	sort.Strings(repo.Topics)

	log.Printf("Comparing: %s with %s", r, *RepositoryToRepo(repo))
	return reflect.DeepEqual(*RepositoryToRepo(repo), r), repo, nil
}

func (r Repo) sync(org GithubOrg, repo *github.Repository) error {

	repo.Description = r.Description
	repo.Homepage = r.Homepage
	repo.Private = r.Private
	repo.HasIssues = r.HasIssues
	repo.HasProjects = r.HasProjects
	repo.HasWiki = r.HasWiki
	repo.Archived = r.Archived
	repo.AllowSquashMerge = r.AllowSquashMerge
	repo.AllowMergeCommit = r.AllowMergeCommit
	repo.AllowRebaseMerge = r.AllowRebaseMerge

	sort.Strings(repo.Topics)
	sort.Strings(r.Topics)

	if !reflect.DeepEqual(repo.Topics, r.Topics) {
		repo.Topics = r.Topics
		err := org.SyncTopics(repo)
		if err != nil {
			return err
		}
	}

	return org.Sync(repo)
}

func (r Repo) Reconcile(org GithubOrg) error {
	repo, err := org.GetRepository(*r.Name)

	log.Printf("Reconciling %s", r)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// if repo is archived remotelly we can skip it, we are not touching archived repos
	if *repo.Archived {
		fmt.Printf("Repo %s is archived remotely, software will ignore archived repos\n", *repo.FullName)
		return nil
	}

	if *r.Archived {
		log.Printf("Local repo %s archived, going to archive remote repo\n", *repo.FullName)
		*repo.Archived = true
		err = org.Sync(repo)
		if err != nil {
			log.Println(err)
		}
	}

	// if repo is not archived we should try to sync the repo
	return r.sync(org, repo)
}

func (org GithubOrg) GetRepository(reponame string) (*github.Repository, error) {

	repo, _, err := org.client.Repositories.Get(org.ctx, org.organization, reponame)

	return repo, err
}

func (org GithubOrg) GetRepositories() ([]*github.Repository, error) {

	// list all repositories for the authenticated user
	repos, _, err := org.client.Repositories.List(org.ctx, org.organization, nil)

	return repos, err
}

func (org GithubOrg) GetRepo(reponame string) (*Repo, error) {
	repository, err := org.GetRepository(reponame)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	repo := RepositoryToRepo(repository)

	return repo, nil
}

func (org GithubOrg) GetRepos() ([]*Repo, error) {
	repos := []*Repo{}
	repositories, err := org.GetRepositories()

	if err != nil {
		return repos, err
	}

	for _, repo := range repositories {
		repo := RepositoryToRepo(repo)

		// Add teams
		repo.Teams, err = org.ListTeams(*repo.Name)
		if err != nil {
			return repos, err
		}

		// Add collaborators
		repo.Collaborators, err = org.ListCollaborators(*repo.Name)
		if err != nil {
			return repos, err
		}

		repos = append(repos, repo)
	}

	return repos, err
}

func (org GithubOrg) ListCollaborators(repo string) ([]*Collaborator, error) {
	cs := []*Collaborator{}

	collaborators, _, err := org.client.Repositories.ListCollaborators(org.ctx, org.organization, repo, &github.ListCollaboratorsOptions{})

	if err != nil {
		return cs, err
	}

	for _, collaborator := range collaborators {
		cs = append(cs, &Collaborator{
			Login:       collaborator.Login,
			Name:        collaborator.Name,
			Permissions: *collaborator.Permissions,
		})
	}

	return cs, nil
}

func (org GithubOrg) ListTeams(repo string) ([]*Team, error) {
	ts := []*Team{}

	teams, _, err := org.client.Repositories.ListTeams(org.ctx, org.organization, repo, &github.ListOptions{})
	if err != nil {
		return ts, err
	}

	for _, team := range teams {
		ts = append(ts, &Team{
			Name:       team.Name,
			Slug:       team.Slug,
			ID:         team.ID,
			Permission: team.Permission,
		})
	}

	return ts, nil
}
