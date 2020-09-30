---
title: Changelog
chapter: true
---
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.15.1] - 2020-09-30

### Fixed

- if the workspace has more the one page of projects, in interactive mode, only the first page was
  being shown. now fixed to run over all pages to fill the list.

### Added

- "Getting Started" section on README.md to help new users to setup theirs environment.

## [v0.15.0] - 2020-09-12

### Added

- support for command line completion on `fish`, `bash` and `zsh` for subcommands and flag name's
- command line completion for arguments and flags for Tags, Projects, Workspaces and Users.
- alias `remove` to command `delete`

### Changed

- using the API `v1` version to get tags available to a workspace.
- `api.Client.Workspaces` renamed to `api.Client.GetWorkspaces` to follow pattern used on other
  functions.
- command `config`, `config set` and `config init` combined to be only one command `config`
- improvements on help of many commands to show usable values.
- `github.com/spf13/cobra` updated to latest possible current version to use completion improvements
  not yet released
- "interactive mode" functions moved to a separate package.

## [v0.14.1] - 2020-09-09

### Fixed

- the project select on interactive mode was not respecting the "default" project when cloning
  or informed through flags/parameters

## [v0.14.0] - 2020-09-08

### Changed

- ask for "interactive mode" and "auto-closing" global configurations on `config init` command.

## [v0.13.0] - 2020-09-08

### Added

- select and multi-select interactive now support "glob like" expressions to filter a option

### Changed

- client name of a project is shown on interactive mode to help identify the project.

### Fixed

- select and multi-select options now support "non-english" characters like "á" by converting then to a ASCII equivalent character.

## [v0.12.2] - 2020-09-04

### Fixed

- flag `--token` help was not showing the right env var name.

## [v0.12.1] - 2020-08-22

### Added

- "How to install" section on README to help new users to understand which options are available.

### Fixed

- improving the "homebrew tap" to allow installation using: `brew install lucassabreu/tap/clockify-cli`

## [v0.12.0] - 2020-08-31

### Added

- support to homebrew for macOs users.

## [v0.11.0] - 2020-08-22

### Added

- new `delete` command to remove a existing time entry from a workspace.
- `edit` command support to interactive mode.

### Fixed

- when cloning a time entry, using interactive mode, the tags selected were not being respected.
- `edit` command was removing all data from time-entry if the flag to fill the field was not being set.

## [v0.10.1] - 2020-08-10

### Fixed

- `in` and `manual` command were showing a error "Project '' informed was not found", even
  when no project id/name is informed, this is now fixed.

## [v0.10.0] - 2020-08-07

### Added

- `clone` command now allow to change the project and description on the
  time entry creation, interactive mode already had this possibility
- new flag `archived` on `project list` to list archived projects
- a new global config `allow-project-name` that, when enabled, allow the user to the project
  name (or parts of it) to be used where the project id is asked.
- common function to get all pages on a paginated request, to not reimplement it, and guarantee
  all entities are being used/returned.

### Fixed

- `clone` sub-command was not asking to confirm the tags when the original time entry already
  had some.
- `clone` command now will respect flags `--tags` and `--when-to-close`.
- "billable" attribute was not being cloned
- keep the current CHANGELOG when extracting the last tag
- some grammatic errors ("applyied" => applied)
- remove mentions to GitHub or Trello token, until integration is implemented

## [v0.9.0] - 2020-07-20

### Added

- new sub-command `version` to allow a quick way to know which version is installed
- sub-command `report` now supports `this-week` and `last-week` as time range aliases
  listing respectively all entries which start this week, and all entries that happened
  on previous week.

### Changed

- all relevant errors now have stack trace with then, which will be printed when the
  flag `--debug` is used.
- error reporting now centralized, removing the need for a helper function in each
  sub-command
- `report`command default output (table) with show in which day the times entries were made.

## [v0.8.1] - 2020-07-09

### Fixed

- `clone` sub-command was not working because the `no-closing` viper config was being
  connected with a non-existing `--no-closing` flag in the `in` sub-command, that does
  not exist anymore.

## [v0.8.0] - 2020-07-08

### Added

- created a new sub-command `manual` that will allow to create "completed" time entries
  in a more easy way.
- created a new flag `--when-to-close` on `in` and `clone` to set close time for the
  time entry being started (if wanted).

### Changed

- `clone` sub-command allows the flag `--no-closing` and will have the same flags as
  `in` to set start and end time (if wanted)
- `in` sub-command will always stops time entries that are open in the moment of the
  sub-command call.
- some helps and messages were improved to better describe what the command does

### Removed

- flags `--trello-token` and `--github-token` were removed because they are not
  currently used and may give false impressions about the cli

### Fixed

- some code for the in and clone sub-commands were duplicated, now they are in `newEntry`
  function that they both used.

## [v0.7.2] - 2020-06-21

### Fixed

- using JSON to notify Netlify, to prevent "malformed url errors"

## [v0.7.1] - 2020-06-21

### Fixed

- snapcraft build/release problems after Travis config update

## [v0.7.0] - 2020-06-21

### Added

- build every pull request as a snapshot to check if it is not failing
- command to auto-generated hugo formatted markdown files from the commands
- implemented a site to better help people to understand what the CLI does,
  without having to download it (live on: https://clockify-cli.netlify.app/)

### Changed

- improved headers on the CHANGELOG to better represent the hierarchies
- moved `in clone` to be just `clone`

### Fixed

- missing release links for the title on the CHANGELOG
- filling the brackets on the LICENSE file

## [v0.6.1] - 2020-06-16

### Added

- `config` command can print the "global" parameters in `json` or `yaml`
- `config` now accepts a argument, which is the name of the parameter,
  when informed only this parameter will be printed

## [v0.6.0] - 2020-06-16

### Added

- some badges, who does not like they?

### Fixed

- help was showing `CLOCKIFY_WROKSPACE` as env var for workspace, the right name is
  `CLOCKIFY_WORKSPACE`
- fixed some `golint` warnings

### Changed

- go mod dependencies updated
- `snapcraft` package only requires network

### Removed

- Removed `GetCurrentUser` in favor of `GetMe` to be closer to the APIs format

## [v0.5.0] - 2020-06-15

### Changed

- `in`, `log` and `report` now don't require you to inform a "user-id", if none is set,
  than will get the user id from the token used to access the api

### Added

- `me` command returns information about the user who owns the token used to access
  the clockify's api

## [v0.4.0] - 2020-06-01

### Added

- table format will show time entry tags

### Changed

- when adding fake entries with `--fill-missing-dates`, will set end time as equal
  to start time, so the duration will be 0 seconds

## [v0.3.2] - 2020-05-22

### Changed

- printing duration as "h:mm:ss" instead of the Go's default format,
  because is more user and sheet applications friendly.

## [v0.3.1] - 2020-04-01

### Fixed

- fixed `--no-closing` being ignored
- interactive flow of `clone` was keeping previous time interval

## [v0.3.0] - 2020-04-01

### Fixed

- minor grammar bug fixes

### Changed

- improvements to the code moving interactive logic of the "in" command into `cmd/common.go`
- "in clone" is now interactive and will ask the user to confirm the time entry data before
  creating it.

## [v0.2.2] - 2020-03-18

### Fixed

- the endpoint `workspaces/<workspace-id>/tags/<tag-id>` does not exist anymore, instead the
  `api.Client` will get all tags of the workspace (`api.Client.GetTags`) and filter the response
  to find the tag by its id.

## [v0.2.1] - 2020-03-02

### Fixed

- `clockify-cli report` parameter `--fill-missing-dates`, was not working

## [v0.2.0] - 2020-03-02

### Added

- `clockify-cli report --fill-missing-dates` when this parameters is set, if there
  are dates from the range informed, will be created "stub" entries to better show
  that are missing entries.

## [v0.1.7] - 2020-02-03

### Added

- `api.Client` now supports getting one specific time entry from a workspace,
  without the need to paginate through all time entries to find it (`GetTimeEntry`
  function).

### Fixed

- `clockify-cli report` was not getting all pages from the period, implemented
  support for pagination and to get "all pages" at once into `Client.Log` and
  `Client.LogRange`

### Changed

- updated README, so it shows the `--help` output as it is now

## [v0.1.6] - 2020-02-03

### Fixed

- fixed bug after Clockify's API changed, where `user` and `project` are not
  automatically provided by the "time-entries" endpoint, unless sending
  an extra parameter `hydrated=true`, and `user` is not provided anymore, so
  now we find it using the user id from the function filter

## [v0.1.5] - 2020-01-08

### Fixed
- fixed bug on the `log` commands, where the previews api url is not available
  anymore, now using `v1/workspace/{workspace}/user/{user}/times-entries`
- spelling of some words fixed and improving some aspects of the code

### Changed
- `go.mod` updated

### Added
- seamless support for query parameters using the interface `QueryAppender`
- support for retrieving the current user of the token (`v1/user`) in the API client.
- `.nvimrc` added to provide spell check

## [v0.1.4] - 2019-08-05

### Added
- Permissions to `snap` installation, so configuration file can be used

## [v0.1.3] - 2019-08-02

### Changed
- Set `publish` to `true` so it will be sent to `snapcraft`

## [v0.1.2] - 2019-08-02

### Added
- Add release to snapcraft by the name `clockify-cli`
- Add command `clockify-cli report` implemented to generate bigger exports. CSV, JSON,
`gofmt` and table formats allowed in this command.

## [v0.1.1] - 2019-06-10

### Changed
- The list returned by the `log` command will the sorted starting from the oldest
time entry.

## [v0.1.0] - 2019-04-08

### Added
- Add `goreleaser` to manage binary and releases of the command
- `clockify-cli in` asks user about new entry information when `interactive` is
  enabled
- Command `clockify-cli config init` allows to start a fresh setup, creating a
  configuration file
- Command `clockify-cli config set` updates/creates one configuration key into the
  configuration file
- `clockify-cli in` commands now allow more flexible time format inputs, can be:
  hh:mm, hh:mm:ss, yyyy-mm-dd hh:mm or yyyy-mm-dd hh:mm:ss
- Command `clockify-cli out` implemented, it will close any pending time entry,
  and show the last entry info when closing it with success
- Command `clockify-cli in clone` implemented, to allow creation of new time
  entries based on existing ones, it also close pending ones, if any
- Command `clockify-cli project list` was implemented, it allows to list the
  projects of a workspace, format the return to table, json, and just id. Helps
  with script automation
- Using https://github.com/spf13/viper to link environment variables and configuration
  files with the global flags. User can set variables `CLOCKIFY_TOKEN`,
  `CLOCKIFY_WORKSPACE` and `CLOCKIFY_USER_ID` instead of using the command flags
- Command `clockify-cli tags` created, to list workspace tags
- Command `clockify-cli in` implemented, to allow creation of new time entries,
  it also close pending ones, if any
- Command `clockify-cli edit <id>` implemented, to allow updates on time entries,
  including the in-progress one using the id: "current
- `--debug` option to allow better understanding of the requests
- Command `clockify-cli log in-progress` implemented, with options to format the
  output, and in the TimeEntry format, instead of TimeEntryImpl
- Command `clockify-cli log` implemented, with options to format the  output,
  will require the user for now
- Package `dto` created to hold all payload objects
- Package `api.Client` to call Clockfy's API
- Command `clockify-cli workspaces` created, with options to format the output
- Command `clockify-cli workspaces users` created, with options to format the
  output to allow retrieving the user's ID

## [v0.0.1] - 2019-03-03
### Added
- This CHANGELOG file to hopefully serve as an evolving example of a
  standardized open source project CHANGELOG.
- README now show which features are expected, and that nothings is done yet
- Golang CLI using [cobra](https://github.com/spf13/cobra)
- Makefile to help setup actions

[Unreleased]: https://github.com/lucassabreu/clockify-cli/compare/v0.15.1...HEAD
[v0.15.1]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.15.1
[v0.15.0]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.15.0
[v0.14.1]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.14.1
[v0.14.0]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.14.0
[v0.13.0]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.13.0
[v0.12.1]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.12.1
[v0.12.0]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.12.0
[v0.11.0]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.11.0
[v0.10.1]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.10.1
[v0.10.0]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.10.0
[v0.9.0]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.9.0
[v0.8.1]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.8.1
[v0.8.0]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.8.0
[v0.7.2]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.7.2
[v0.7.1]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.7.1
[v0.7.0]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.7.0
[v0.6.1]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.6.1
[v0.6.0]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.6.0
[v0.5.0]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.5.0
[v0.4.0]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.4.0
[v0.3.2]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.3.2
[v0.3.1]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.3.1
[v0.3.0]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.3.0
[v0.2.2]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.2.2
[v0.2.1]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.2.1
[v0.2.0]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.2.0
[v0.1.7]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.1.7
[v0.1.6]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.1.6
[v0.1.5]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.1.5
[v0.1.4]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.1.4
[v0.1.3]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.1.3
[v0.1.2]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.1.2
[v0.1.1]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.1.1
[v0.1.0]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.1.0
[v0.0.1]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.0.1
