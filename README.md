clockify-cli
============

A simple cli to manage your time entries on Clockify from terminal

Features
--------

* [X] List time entries from a day
  * [X] List in progress entry
* [X] Start a new time entry
  + [ ] Using a GitHub issue
  + [ ] Using a Trello card
  + [X] Cloning last time entry
* [ ] Stop the last entry
* [X] List workspace projects
* [ ] Link a Clockify Project with Github:Issues repository
* [ ] Link a Clockify Project with Trello board
* [X] List Clockify Workspaces
* [X] List Clockify Workspaces Users
* [X] List Clockify Tags
* [X] Edit time entry

Help
----

```
Allow to integrate with Clockify through terminal

Usage:
  clockify-cli [command]

Available Commands:
  edit        Edit a time entry, use id "current" to apply to time entry in progress
  help        Help about any command
  in          Create a new time entry and starts it
  log         List the entries from a specific day
  out         Stops the last time entry
  project     Allow project aliasing and integration of a project with GitHub:Issues or Trello
  tags        List tags of workspace
  workspaces  List user's workspaces

Flags:
      --config string         config file (default is $HOME/.clockify-cli.yaml)
      --debug                 show debug log (defaults to env $CLOCKIFY_DEBUG)
      --github-token string   gitHub's token (defaults to env $CLOCKIFY_GITHUB_TOKEN)
  -h, --help                  help for clockify-cli
  -t, --token string          clockify's token (defaults to env $CLOCKIFY_TOKEN)
                              	Can be generated here: https://clockify.me/user/settings#generateApiKeyBtn
      --trello-token string   trello's token (defaults to env $CLOCKIFY_TRELLO_TOKEN)
  -u, --user-id string        user id from the token (defaults to env $CLOCKIFY_USER_ID)
  -w, --workspace string      workspace to be used (defaults to env $CLOCKIFY_WROKSPACE)

Use "clockify-cli [command] --help" for more information about a command.
```
