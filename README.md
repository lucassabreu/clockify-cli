Clockify CLI
============

A simple cli to manage your time entries on Clockify from terminal

[![clockify-cli](https://snapcraft.io//clockify-cli/badge.svg?classes=badges)](https://snapcraft.io/clockify-cli)
[![Build Status](https://travis-ci.org/lucassabreu/clockify-cli.svg?branch=master&classes=badges)](https://travis-ci.org/lucassabreu/clockify-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/lucassabreu/clockify-cli?classes=badges)](https://goreportcard.com/report/github.com/lucassabreu/clockify-cli)
[![Netlify Status](https://api.netlify.com/api/v1/badges/8667b9f6-4ca2-4ee4-865e-20b5848e7059/deploy-status?classes=badges)](https://app.netlify.com/sites/clockify-cli/deploys)

Features
--------

* [X] List time entries from a day
  * [X] List in progress entry
* [X] Report time entries using a date range
  * [X] Inform date range as parameters
  * [X] "auto filter" for last month
  * [X] "auto filter" for this month
* [X] Start a new time entry
  + [ ] Using a GitHub issue
  + [ ] Using a Trello card
  + [X] Cloning last time entry
  + [X] Ask input interactively
* [X] Stop the last entry
* [X] List workspace projects
* [ ] Link a Clockify Project with Github:Issues repository
* [ ] Link a Clockify Project with Trello board
* [X] List Clockify Workspaces
* [X] List Clockify Workspaces Users
* [X] List Clockify Tags
* [X] Edit time entry
* [X] Configuration management
  * [X] Initialize configuration
  * [X] Update individual configuration
  * [X] Show current configuration

Help
----

```
Allow to integrate with Clockify through terminal

Usage:
  clockify-cli [command]

Available Commands:
  clone       Copy a time entry and starts it (use "last" to copy the last one)
  config      Manages configuration file parameters
  edit        Edit a time entry, use id "current" to apply to time entry in progress
  help        Help about any command
  in          Create a new time entry and starts it
  log         List the entries from a specific day
  me          Show the user info
  out         Stops the last time entry
  project     Allow project aliasing and integration of a project with GitHub:Issues or Trello
  report      report for date ranges and with more data (format date as 2016-01-02)
  tags        List tags of workspace
  workspaces  List user's workspaces

Flags:
      --config string         config file (default is $HOME/.clockify-cli.yaml)
      --debug                 show debug log (defaults to env $CLOCKIFY_DEBUG)
      --github-token string   GitHub's token (defaults to env $CLOCKIFY_GITHUB_TOKEN)
  -h, --help                  help for clockify-cli
  -i, --interactive           show interactive log (defaults to env $CLOCKIFY_INTERACTIVE)
  -t, --token string          clockify's token (defaults to env $CLOCKIFY_TOKEN)
                                Can be generated here: https://clockify.me/user/settings#generateApiKeyBtn
      --trello-token string   Trello's token (defaults to env $CLOCKIFY_TRELLO_TOKEN)
  -u, --user-id string        user id from the token (defaults to env $CLOCKIFY_USER_ID)
  -w, --workspace string      workspace to be used (defaults to env $CLOCKIFY_WORKSPACE)

Use "clockify-cli [command] --help" for more information about a command.
```
