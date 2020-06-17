# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## Changed

- go mod dependencies updated
- `snapcraft` package only requires network

## [v0.5.0] - 2020-06-15

## Changed

- `in`, `log` and `report` now don't require you to inform a "user-id", if none is set,
  than will get the user id from the token used to access the api

## Added

- `me` command returns information about the user who owns the token used to access
  the clockify's api

## [v0.4.0] - 2020-06-01

## Added

- table format will show time entry tags

## Changed

- when adding fake entries with `--fill-missing-dates`, will set end time as equal
  to start time, so the duration will be 0 seconds

## [v0.3.2] - 2020-05-22

## Changed

- printing duration as "h:mm:ss" instead of the Go's default format,
  because is more user and sheet applications friendly.

## [v0.3.1] - 2020-04-01

## Fixed

- fixed `--no-closing` being ignored
- interactive flow of `clone` was keeping previous time interval

## [v0.3.0] - 2020-04-01

## Fixed

- minor grammar bug fixes

## Changed

- improvements to the code moving interactive logic of the "in" command into `cmd/common.go`
- "in clone" is now interactive and will ask the user to confirm the time entry data before
  creating it.

## [v0.2.2] - 2020-03-18

## Fixed

- the endpoint `workspaces/<workspace-id>/tags/<tag-id>` does not exist anymore, instead the
  `api.Client` will get all tags of the workspace (`api.Client.GetTags`) and filter the response
  to find the tag by its id.

## [v0.2.1] - 2020-03-02

## Fixed

- `clockify-cli report` parameter `--fill-missing-dates`, was not working

## [v0.2.0] - 2020-03-02

## Added

- `clockify-cli report --fill-missing-dates` when this parameters is set, if there
  are dates from the range informed, will be created "stub" entries to better show
  that are missing entries.

## [v0.1.7] - 2020-02-03

## Added

- `api.Client` now supports getting one specific time entry from a workspace,
  without the need to paginate through all time entries to find it (`GetTimeEntry`
  function).

## Fixed

- `clockify-cli report` was not getting all pages from the period, implemented
  support for pagination and to get "all pages" at once into `Client.Log` and
  `Client.LogRange`

## Changed

- updated README, so it shows the `--help` output as it is now

## [v0.1.6] - 2020-02-03

## Fixed

- fixed bug after Clockify's API changed, where `user` and `project` are not
  automatically provided by the "time-entries" endpoint, unless sending
  an extra parameter `hydrated=true`, and `user` is not provided anymore, so
  now we find it using the user id from the function filter

## [v0.1.5] - 2020-01-08

## Fixed
- fixed bug on the `log` commands, where the previews api url is not available
  anymore, now using `v1/workspace/{workspace}/user/{user}/times-entries`
- spelling of some words fixed and improving some aspects of the code

## Changed
- `go.mod` updated

## Added
- seamless support for query parameters using the interface `QueryAppender`
- support for retrieving the current user of the token (`v1/user`) in the API client.
- `.nvimrc` added to provide spell check

## [v0.1.4] - 2019-08-05

## Added
- Permissions to `snap` installation, so configuration file can be used

## [v0.1.3] - 2019-08-02

## Changed
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

[Unreleased]: https://github.com/lucassabreu/clockify-cli/compare/v0.1.4...HEAD
[v0.1.4]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.1.4
[v0.1.3]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.1.3
[v0.1.2]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.1.2
[v0.1.1]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.1.1
[v0.1.0]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.1.0
[v0.0.1]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.0.1
