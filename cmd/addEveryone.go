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

	"github.com/spf13/cobra"

	"log"

	"github.com/getsentry/raven-go"
	"github.com/google/go-github/v28/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// addEveryoneCmd represents the add-everyone command
var addEveryoneCmd = &cobra.Command{
	Use:   "add-everyone",
	Short: "Assign every member of the organization to a target team.",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		accessToken := viper.Get("accessToken").(string)
		targetOrg := cmd.Flag("org").Value.String()
		targetTeam := cmd.Flag("team").Value.String()

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)

		teams, _, err := client.Teams.ListTeams(ctx, targetOrg, nil)
		if err != nil {
			log.Printf("Team `%s` is not found in the organization `%s`!", targetTeam, targetOrg)
			raven.CaptureErrorAndWait(err, nil)
			log.Fatal(err)
		}
		team := findTeam(teams, targetTeam)
		if team == nil {
			newTeamPrivacy := "closed"
			newTeam := github.NewTeam{Name: targetTeam, Privacy: &newTeamPrivacy}
			team, _, err = client.Teams.CreateTeam(ctx, targetOrg, newTeam)
			if err != nil {
				log.Println(err)
				raven.CaptureErrorAndWait(err, nil)
				log.Fatalf("Failed to create the new team `%s`  in the organization `%s`!", targetTeam, targetOrg)
			}
		}

		userLogins := getUserLogins(client, ctx, targetOrg, *team.ID)
		opt := &github.TeamAddTeamMembershipOptions{Role: "member"}

		for _, userLogin := range userLogins.ToSlice() {
			login := userLogin.(string)
			_, _, err := client.Teams.AddTeamMembership(ctx, *team.ID, login, opt)
			if err != nil {
				log.Printf("Failed to add a user `%s` to the team `%s`: ", login, targetTeam)
				log.Print(err)
			}
			log.Printf("`%s` is now a member of the team `%s`.", login, targetTeam)
		}

		log.Println("Done!")
	},
}

func init() {
	RootCmd.AddCommand(addEveryoneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addEveryoneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addEveryoneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	addEveryoneCmd.Flags().String("org", "", "GitHub organization")
	addEveryoneCmd.MarkFlagRequired("org")
	addEveryoneCmd.Flags().String("team", "", "Team which every member of the organization belongs to")
	addEveryoneCmd.MarkFlagRequired("team")
}
