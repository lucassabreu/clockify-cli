# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Add comand `clockify-cli report` implemented to generate bigger exports. CSV, JSON,
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
- Command `clockify-cli config set` updates/creates one config key into the
  config file
- `clockify-cli in` commands now allow more flexible time format inputs, can be:
  hh:mm, hh:mm:ss, yyyy-mm-dd hh:mm or yyyy-mm-dd hh:mm:ss
- Command `clockify-cli out` implemented, it will close any pending time entry,
  and show the last entry info when closing it with success
- Command `clockify-cli in clone` implemented, to allow creation of new time
  entries based on existing ones, it also close pending ones, if any
- Command `clockify-cli project list` was implemented, it allows to list the
  projects of a worspace, format the return to table, json, and just id. Helps
  with script automatation
- Using https://github.com/spf13/viper to link enviroment variables and config
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
  output to allow retriving the user's ID

## [0.0.1] - 2019-03-03
### Added
- This CHANGELOG file to hopefully serve as an evolving example of a
  standardized open source project CHANGELOG.
- README now show which features are expected, and that nothings is done yet
- Golang CLI using [cobra](https://github.com/spf13/cobra)
- Makefile to help setup actions

[Unreleased]: https://github.com/lucassabreu/clockify-cli/compare/v0.0.1...HEAD
[0.0.1]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.0.1
