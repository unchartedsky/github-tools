// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"github.com/deckarep/golang-set"

	"github.com/spf13/viper"
	"log"
	"os"
)

// addEveryoneCmd represents the add-everyone command
var addEveryoneCmd = &cobra.Command{
	Use:   "add-everyone",
	Short: "Assign every member of the organization to a target team.",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		accessToken := viper.Get("accessToken").(string)
		if len(accessToken) ==0 {
			log.Fatal("GitHub access token is required!")
		}

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)


		teams, _, err := client.Organizations.ListTeams(ctx, Org, nil)
		if err != nil {
			log.Printf("Team `%s` is not found in the organization `%s`!", Team, Org)
			log.Fatal(err)
		}
		team := Find(teams, Team)
		if team == nil {
			log.Fatalf("Team `%s` is not found in the organization `%s`!", Team, Org)
		}

		orgUsers, _, err := client.Organizations.ListMembers(ctx, Org, nil)
		teamUsers, _, err := client.Organizations.ListTeamMembers(ctx, *team.ID, nil)

		numberOfNewUsers := len(orgUsers) - len(teamUsers)
		if numberOfNewUsers == 0 {
			log.Printf("No new users are found.")
			os.Exit(0)
		}
		log.Printf("%d new users are found", len(orgUsers)-len(teamUsers))

		userLogins := NewUserLogins(orgUsers, teamUsers)

		for _, userLogin := range userLogins.ToSlice() {
			login := userLogin.(string)
			_, _, err := client.Organizations.AddTeamMembership(ctx, *team.ID, login, nil)
			if err != nil {
				log.Printf("Failed to add a user `%s` to the team `%s`: ", login, Team)
				log.Print(err)
			}
			log.Printf("`%s` is now a member of the team `%s`.", login, Team)
		}

		fmt.Println("Done!")
	},
}

var Team string
var Org string

func init() {
	RootCmd.AddCommand(addEveryoneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// add-everyoneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// add-everyoneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	addEveryoneCmd.Flags().StringVar(&Org, "org", "" , "GitHub organization")
	addEveryoneCmd.MarkFlagRequired("org")
	addEveryoneCmd.Flags().StringVar(&Team, "team", "", "Team which every member of the organization belongs to")
	addEveryoneCmd.MarkFlagRequired("team")

}

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
	orgLogins := mapset.NewSetFromSlice( orgUserIds)

	teamUserIds := Map(teamUsers)
	teamLogins:= mapset.NewSetFromSlice(teamUserIds)

	newUserLogins := orgLogins.Difference(teamLogins)

	return newUserLogins
}