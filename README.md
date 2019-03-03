clockify-cli
============

A simple cli to manage your time entries on Clockify from terminal

Features
--------

* [X] List time entries from a day
* [ ] Start a new time entry
  + [ ] Using a GitHub issue
  + [ ] Using a Trello card
* [ ] Stop the last entry
* [ ] List workspace projects
* [ ] Link a Clockify Project with Github:Issues repository
* [ ] Link a Clockify Project with Trello board
* [X] List Clockify Workspaces
* [X] List Clockify Workspaces Users

Help
----

```
Allow to integrate with Clockify through terminal

Usage:
  clockify-cli [command]

Available Commands:
  help        Help about any command
  in          Create a new time entry and starts it
  log         List the entries from a specific day
  out         Stops the last time entry
  project     Allow project aliasing and integration of a project with GitHub:Issues or Trello
  workspaces  List user's workspaces

Flags:
      --config string         config file (default is $HOME/.clockify-cli.yaml)
      --debug                 show debug log
      --github-token string   gitHub's token
  -h, --help                  help for clockify-cli
  -t, --token string          clockify's token, can be generated here: https://clockify.me/user/settings#generateApiKeyBtn
      --trello-token string   trello's token
  -u, --user-id string        user id from the token
  -w, --workspace string      workspace to be used

Use "clockify-cli [command] --help" for more information about a command.
```
