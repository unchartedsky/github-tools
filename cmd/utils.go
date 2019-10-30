package cmd

import (
	"context"
	"log"
	"os"
	"sync"

	mapset "github.com/deckarep/golang-set"
	"github.com/getsentry/raven-go"
	"github.com/google/go-github/v28/github"
)

func findTeam(teams []*github.Team, teamName string) *github.Team {
	for _, team := range teams {
		if teamName == *team.Name {
			return team
		}
	}
	return nil
}

func userToId(vs []*github.User) []interface{} {
	vsm := make([]interface{}, len(vs))
	for i, v := range vs {
		vsm[i] = *v.Login
	}
	return vsm
}

func getUserLogins(client *github.Client, ctx context.Context, org string, team int64) mapset.Set {
	var wg sync.WaitGroup
	wg.Add(2)

	var orgUsers []*github.User
	go func() {
		opt := &github.ListMembersOptions{ListOptions: github.ListOptions{PerPage: 30}}
		for {
			users, resp, err := client.Organizations.ListMembers(ctx, org, opt)
			if err != nil {
				raven.CaptureErrorAndWait(err, nil)
				log.Fatal(err)
			}

			orgUsers = append(orgUsers, users...)

			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}

		wg.Done()
	}()

	var teamUsers []*github.User
	go func() {
		opt := &github.TeamListTeamMembersOptions{ListOptions: github.ListOptions{PerPage: 30}}
		for {
			users, resp, err := client.Teams.ListTeamMembers(ctx, team, opt)
			if err != nil {
				raven.CaptureErrorAndWait(err, nil)
				log.Fatal(err)
			}

			teamUsers = append(teamUsers, users...)

			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}
		wg.Done()
	}()

	wg.Wait()

	if orgUsers == nil {
		log.Fatal("Getting the list of organization members failed!")
	}
	if teamUsers == nil {
		log.Fatal("Getting the list of team members failed!")
	}

	numberOfNewUsers := len(orgUsers) - len(teamUsers)
	if numberOfNewUsers == 0 {
		log.Printf("No new users are found.")
		os.Exit(0)
	}
	log.Printf("%d new users are found", len(orgUsers)-len(teamUsers))
	return newUserLogins(orgUsers, teamUsers)
}

func newUserLogins(orgUsers []*github.User, teamUsers []*github.User) mapset.Set {
	// TODO 깔끔한 알고리즘으로 나중에 교체하자
	orgUserIds := userToId(orgUsers)
	orgLogins := mapset.NewSetFromSlice(orgUserIds)

	teamUserIds := userToId(teamUsers)
	teamLogins := mapset.NewSetFromSlice(teamUserIds)

	newUserLogins := orgLogins.Difference(teamLogins)

	return newUserLogins
}

func contains(items []string, search string) bool {
	for _, value := range items {
		if value == search {
			return true
		}
	}
	return false
}
