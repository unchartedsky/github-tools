[![Build Status](https://travis-ci.org/UnchartedSky/github-tools.svg?branch=master)](https://travis-ci.org/UnchartedSky/github-tools)
[![Docker Pulls](https://img.shields.io/docker/pulls/mashape/kong.svg)](https://hub.docker.com/r/unchartedsky/github-tools)

# GitHub Tools

``` bash
$ cat ~/.github-tools.yaml 
accessToken: "YOUR_GITHUB_ACCESS_TOKEN"
```

``` bash
$ go run main.go 
GitHub management tool

Usage:
  github-tools [command]

Available Commands:
  add-everyone Assign every member of the organization to a target team.
  add-team     Add a team to all the repositories, which belong to the organization
  help         Help about any command

Flags:
      --config string   config file (default is $HOME/.github-tools.yaml)
  -h, --help            help for github-tools
  -t, --toggle          Help message for toggle
      --token string    GitHub access token

Use "github-tools [command] --help" for more information about a command.

```

``` bash
go run main.go add-everyone --team Everyone --org UnchartedSky
```

``` bash
go run main.go add-team --team Everyone --org UnchartedSky
```

## Integrations with Sentry.io

Specify the `SENTRY_DSN` environment variable to use Sentry.io as your monitoring tool.