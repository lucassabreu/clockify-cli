# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
