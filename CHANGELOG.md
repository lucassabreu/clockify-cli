# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Using https://github.com/spf13/viper to link enviroment variables and config 
  files with the global flags. User can set variables `CLOCKIFY_TOKEN`,
  `CLOCKIFY_WORKSPACE` and `CLOCKIFY_USER_ID` instead of using the command flags
- `--debug` option to allow better understanding of the requests
- Command `clockify-cli log in-progress` implemented, with options to format the
  output, and in the TimeEntry format, instead of TimeEntryImpl
- Command `clockify-cli log` implemented, with options to format the 
  output, will require the user for now
- Package `dto` created to hold all payload objects
- Package `api.Client` to call Clockfy's API
- Command `clockify-cli workspaces` created, with options to format the output
- Command `clockify-cli workspaces users` created, with options to format the output
  to allow retriving the user's ID

## [0.0.1] - 2019-03-03
### Added
- This CHANGELOG file to hopefully serve as an evolving example of a
  standardized open source project CHANGELOG.
- README now show which features are expected, and that nothings is done yet
- Golang CLI using [cobra](https://github.com/spf13/cobra)
- Makefile to help setup actions

[Unreleased]: https://github.com/lucassabreu/clockify-cli/compare/v0.0.1...HEAD
[0.0.1]: https://github.com/lucassabreu/clockify-cli/releases/tag/v0.0.1
