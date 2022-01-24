package githuborg

import (
	"fmt"

	"github.com/google/go-github/github"
)

func remote(org GithubOrg) ([]*github.Team, error) {
	teams, _, err := org.client.Teams.ListTeams(org.ctx, org.organization, &github.ListOptions{})

	if err != nil {
		fmt.Println(err)
		return teams, err
	}

	return teams, err
}
