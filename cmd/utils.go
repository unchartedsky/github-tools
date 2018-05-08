package cmd

import (
	"github.com/google/go-github/github"
	"github.com/deckarep/golang-set"
)

func Find(teams []*github.Team, teamName string) *github.Team {
	for _, team := range teams {
		if teamName == *team.Name {
			return team
		}
	}
	return nil
}

func Map(vs []*github.User) []interface{} {
	vsm := make([]interface{}, len(vs))
	for i, v := range vs {
		vsm[i] = *v.Login
	}
	return vsm
}

func NewUserLogins(orgUsers []*github.User, teamUsers []*github.User) mapset.Set {
	// TODO 깔끔한 알고리즘으로 나중에 교체하자
	orgUserIds := Map(orgUsers)
	orgLogins := mapset.NewSetFromSlice(orgUserIds)

	teamUserIds := Map(teamUsers)
	teamLogins:= mapset.NewSetFromSlice(teamUserIds)

	newUserLogins := orgLogins.Difference(teamLogins)

	return newUserLogins
}

func Contains(items []string, search string) bool {
	for _, value := range items {
		if value == search {
			return true
		}
	}
	return false
}