Clockify CLI
============

A simple cli to manage your time entries on Clockify from terminal

[![Release](https://img.shields.io/github/release/lucassabreu/clockify-cli.svg?classes=badges)](https://github.com/lucassabreu/clockify-cli/releases/latest)
[![clockify-cli](https://snapcraft.io//clockify-cli/badge.svg?classes=badges)](https://snapcraft.io/clockify-cli)
[![Build Status](https://github.com/lucassabreu/clockify-cli/actions/workflows/release.yml/badge.svg?classes=badges)](.github/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/lucassabreu/clockify-cli?classes=badges)](https://goreportcard.com/report/github.com/lucassabreu/clockify-cli)
[![Netlify Status](https://api.netlify.com/api/v1/badges/8667b9f6-4ca2-4ee4-865e-20b5848e7059/deploy-status?classes=badges)](https://app.netlify.com/sites/clockify-cli/deploys)

Features
--------

* [x] List time entries from a day
  * [x] List in progress entry
* [x] Report time entries using a date range
  * [x] Inform date range as parameters
  * [x] "auto filter" for last month
  * [x] "auto filter" for this month
* [x] Start a new time entry
  * [ ] Using a GitHub issue
  * [ ] Using a Trello card
  * [x] Cloning last time entry
  * [x] Ask input interactively
* [x] Stop the last entry
* [x] List workspace projects
* [ ] Link a Clockify Project with Github:Issues repository
* [ ] Link a Clockify Project with Trello board
* [x] List Clockify Workspaces
* [x] List Clockify Workspaces Users
* [x] List Clockify Tags
* [x] Edit time entry
* [x] Configuration management
  * [x] Initialize configuration
  * [x] Update individual configuration
  * [x] Show current configuration

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
---------------

After you install the CLI, run `clockify-cli config --init` to setup your environment variables. You’ll be prompted to enter your user information. You can get your clockify api token [here](https://clockify.me/user/settings).

```console
foo@bar:~$ clockify-cli config --init
? User Generated Token: <your-api-token>
? Choose default Workspace: <workspace-id> - John Doe's workspace
? Choose your user: <user-id> - John Doe
? Should try to find project by its name? Yes
? Should use "Interactive Mode" by default? Yes
? Which days of the week do you work? monday, tuesday, wednesday, thursday, friday
? Should allow starting time entries with incomplete data? No
```

The CLI saves your workspace info and an API token to `~/.clockify-cli.yaml` for future use.

> :exclamation: If you have installed the client using `snap` this file will not be accessible to you, but the configs will still be persisted.

Now you’re ready to create your first Clockify entry:

```console
foo@bar:~$ clockify-cli in -i
? Choose your project: project-id - Example Project
? Description: Clockify CLI Test
? Start: (2020-09-04 16:02:45)
? End (leave it blank for empty):
+---------------------+----------+----------+---------+-------------------------------+-------------------+------+
|         ID          |  START   |   END    |   DUR   |            PROJECT            |    DESCRIPTION    | TAGS |
+---------------------+----------+----------+---------+-------------------------------+-------------------+------+
|    time-entry-id    | 16:02:45 | 16:03:47 | 0:01:02 |  project-id - Example Project | Clockify CLI Test |      |
+---------------------+----------+----------+---------+-------------------------------+-------------------+------+
```

After finishing your work you can stop the entry using `clockify-cli out`

```console
foo@bar:~$ clockify-cli out
+---------------------+----------+----------+---------+---------+-------------------+------+
|         ID          |  START   |   END    |   DUR   | PROJECT |    DESCRIPTION    | TAGS |
+---------------------+----------+----------+---------+---------+-------------------+------+
|    time-entry-id    | 16:02:45 | 16:08:06 | 0:05:21 |         | Clockify CLI Test |      |
+---------------------+----------+----------+---------+---------+-------------------+------+
```

If you want to re-start the last entry you can use `clockify-cli clone last`

```console
foo@bar:~$ clockify-cli clone last
+---------------------+----------+----------+---------+---------+-------------------+------+
|         ID          |  START   |   END    |   DUR   | PROJECT |    DESCRIPTION    | TAGS |
+---------------------+----------+----------+---------+---------+-------------------+------+
|    time-entry-id    | 16:10:57 | 16:11:09 | 0:00:12 |         | Clockify CLI Test |      |
+---------------------+----------+----------+---------+---------+-------------------+------+

```

Help
----

```
Allow to integrate with Clockify through terminal

Usage:
  clockify-cli [command]

Available Commands:
  clone         Copy a time entry and starts it (use "last" to copy the last one)
  completion    Generate completion script
  config        Manages configuration file parameters
  delete        Delete time entry(ies), use id "current" to apply to time entry in progress
  edit          Edit a time entry, use id "current" to apply to time entry in progress
  edit-multiple Edit multiple time entries at once, use id "current"/"last" to apply to time entry in progress.
  gendocs       Generate Markdown documentation for the clockify-cli.
  help          Help about any command
  in            Create a new time entry and starts it (will close time entries not closed)
  log           List the entries from a specific day
  manual        Creates a new completed time entry (does not stop on-going time entries)
  me            Show the user info
  out           Stops the last time entry
  project       List projects from a workspace
  report        List all time entries in the date ranges and with more data (format date as 2016-01-02)
  tags          List tags of workspace
  version       Version of the command
  workspace     List user's workspaces

Flags:
      --allow-project-name   allow use of project name when id is asked (defaults to env $CLOCKIFY_ALLOW_PROJECT_NAME)
      --config string        config file (default is $HOME/.clockify-cli.yaml)
      --debug                show debug log (defaults to env $CLOCKIFY_DEBUG)
  -h, --help                 help for clockify-cli
  -i, --interactive          will prompt you to confirm/complement commands input before executing the action (defaults to env $CLOCKIFY_INTERACTIVE).
                             	You can be disable it temporally by setting it to 0 (-i=0 or CLOCKIFY_INTERACTIVE=0)
  -t, --token string         clockify's token (defaults to env $CLOCKIFY_TOKEN)
                             	Can be generated here: https://clockify.me/user/settings#generateApiKeyBtn
  -u, --user-id string       user id from the token (defaults to env $CLOCKIFY_USER_ID)
  -w, --workspace string     workspace to be used (defaults to env $CLOCKIFY_WORKSPACE)

Use "clockify-cli [command] --help" for more information about a command.
```

See more information about the sub-commands at: https://clockify-cli.netlify.app/en/commands/clockify-cli/
