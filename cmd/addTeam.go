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

	"log"

	"github.com/getsentry/raven-go"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// addTeamCmd represents the addTeam command
var addTeamCmd = &cobra.Command{
	Use:   "add-team",
	Short: "Add a team to all the repositories, which belong to the organization",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		accessToken := viper.Get("accessToken").(string)
		targetOrg := cmd.Flag("org").Value.String()
		targetTeam := cmd.Flag("team").Value.String()
		topicToExclude := cmd.Flag("exclude").Value.String()
		permission := cmd.Flag("permission").Value.String()

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)

		teams, _, err := client.Organizations.ListTeams(ctx, targetOrg, nil)
		if err != nil {
			log.Printf("Team `%s` is not found in the organization `%s`!", targetTeam, targetOrg)
			raven.CaptureErrorAndWait(err, nil)
			log.Fatal(err)
		}
		team := findTeam(teams, targetTeam)
		if team == nil {
			raven.CaptureErrorAndWait(err, nil)
			log.Fatalf("Team `%s` is not found in the organization `%s`!", targetTeam, targetOrg)
		}

		opt := &github.RepositoryListByOrgOptions{Type: "all", ListOptions: github.ListOptions{PerPage: 30}}
		for {
			repos, resp, err := client.Repositories.ListByOrg(context.Background(), targetOrg, opt)
			if err != nil {
				raven.CaptureErrorAndWait(err, nil)
				log.Fatal(err)
			}

			for _, repo := range repos {
				repoName := *repo.Name
				//fmt.Println(repoName)

				if contains(repo.Topics, topicToExclude) {
					log.Printf("Repository `%s` has topic `%s` and is skipped.\n", repoName, topicToExclude)
					continue
				}

				option := &github.OrganizationAddTeamRepoOptions{Permission: permission}
				_, err = client.Organizations.AddTeamRepo(ctx, *team.ID, targetOrg, repoName, option)
				if err != nil {
					raven.CaptureErrorAndWait(err, nil)
					log.Fatal(err)
				}
			}

			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}

		log.Println("Done!")
	},
}

func init() {
	RootCmd.AddCommand(addTeamCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addTeamCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addTeamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	addTeamCmd.Flags().String("org", "", "GitHub organization")
	addTeamCmd.MarkFlagRequired("org")
	addTeamCmd.Flags().String("team", "", "Team which every member of the organization belongs to")
	addTeamCmd.MarkFlagRequired("team")
	addTeamCmd.Flags().String("exclude", "private-repository", "The repository with this topic will be excluded")
	addTeamCmd.Flags().String("permission", "push", "Team's permission to the repositories")
}
