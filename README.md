Clockify CLI
============

A simple cli to manage your time entries on Clockify from terminal

[![Release](https://img.shields.io/github/release/lucassabreu/clockify-cli.svg?classes=badges)](https://github.com/lucassabreu/clockify-cli/releases/latest)
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

How to install [![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?classes=badges)](https://github.com/goreleaser)
--------------

#### Using [`homebrew`](https://brew.sh/):

```sh
brew install lucassabreu/tap/clockify-cli
```

#### Using [`snapcraft`](https://snapcraft.io/clockify-cli)

```sh
sudo snap install clockify-cli
```

#### Using `go get`

```sh
go get -u github.com/lucassabreu/clockify-cli
```

#### By Hand

Go to the [releases page](https://github.com/lucassabreu/clockify-cli/releases) and download the pre-compiled
binary that fits your system.

Getting Started
----

After you install the CLI, run the `clockify config init` command to setup your environment variables. You’ll be prompted to enter your user information. You can get your clockify api token [here](https://clockify.me/user/settings).

```console
foo@bar:~$ clockify config init
? User Generated Token: <your-api-token>
? Choose default Workspace: workspace-id - John Doe's workspace
? Choose your user: user-id - John Doe
? Should try to find project by its name? (y/N) Yes
```

The CLI saves your workspace info and an API token to `~/.clockify-cli.yaml` for future use.

Now you’re ready to create your first Clockify entry:

```console
foo@bar:~$ clockify-cli in -i
? Choose your project: project-id - Example Project
? Description: Clockify CLI Test
? Start: (2020-09-04 16:02:45)
? End (leave it blank for empty):
+------------------+----------+----------+---------+---------+-------------------+------+
|        ID        |  START   |   END    |   DUR   | PROJECT |    DESCRIPTION    | TAGS |
+------------------+----------+----------+---------+---------+-------------------+------+
|    project-id    | 16:02:45 | 16:03:47 | 0:01:02 |         | Clockify CLI Test |      |
+------------------+----------+----------+---------+---------+-------------------+------+
```

After finishing your work you can stop the last entry using `clockify-cli out`

```console
foo@bar:~$ clockify-cli out  
+------------------+----------+----------+---------+---------+-------------------+------+
|        ID        |  START   |   END    |   DUR   | PROJECT |    DESCRIPTION    | TAGS |
+------------------+----------+----------+---------+---------+-------------------+------+
|    project-id    | 16:02:45 | 16:08:06 | 0:05:21 |         | Clockify CLI Test |      |
+------------------+----------+----------+---------+---------+-------------------+------+
```

If you want to re-start the last entry in a project you can use `clone last -i`

```console
foo@bar:~$ clone last -i
? Choose your project: project-id - Example Project
? Description: Clockify CLI Test
? Start: 2020-09-04 16:10:57
? End (leave it blank for empty):
+------------------+----------+----------+---------+---------+-------------------+------+
|        ID        |  START   |   END    |   DUR   | PROJECT |    DESCRIPTION    | TAGS |
+------------------+----------+----------+---------+---------+-------------------+------+
|    project-id    | 16:10:57 | 16:11:09 | 0:00:12 |         | Clockify CLI Test |      |
+------------------+----------+----------+---------+---------+-------------------+------+

```

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
  gendocs     Generate Markdown documentation for the clockify-cli.
  help        Help about any command
  in          Create a new time entry and starts it (will close time entries not closed)
  log         List the entries from a specific day
  manual      Creates a new completed time entry (does not stop on-going time entries)
  me          Show the user info
  out         Stops the last time entry
  project     List projects from a workspace
  report      List all time entries in the date ranges and with more data (format date as 2016-01-02)
  tags        List tags of workspace
  workspaces  List user's workspaces

Flags:
      --config string      config file (default is $HOME/.clockify-cli.yaml)
      --debug              show debug log (defaults to env $CLOCKIFY_DEBUG)
  -h, --help               help for clockify-cli
  -i, --interactive        show interactive log (defaults to env $CLOCKIFY_INTERACTIVE)
  -t, --token string       clockify's token (defaults to env $CLOCKIFY_TOKEN)
                           	Can be generated here: https://clockify.me/user/settings#generateApiKeyBtn
  -u, --user-id string     user id from the token (defaults to env $CLOCKIFY_USER_ID)
  -w, --workspace string   workspace to be used (defaults to env $CLOCKIFY_WORKSPACE)

Use "clockify-cli [command] --help" for more information about a command.
```

See more information about the sub-commands at: https://clockify-cli.netlify.app/en/commands/clockify-cli/
