# Contributing

Thank you for the interest in Contributing to Clockify CLI.

We accept pull requests for bug fixes and features (preferably that were discussed on an issue
before). Also opening issues with feature requests and reporting bugs are very important
contributions.

Please do:

- Check in the [issues][issues] if the [bug][bugs] or [feature request][enhancement] has not been submitted.
- Open an issue if things aren't working as expected.
- Open an issue to propose new features or improvements on existing ones.
- Open a pull request to fix a [bug][bugs].
- Open a pull request for any open issue labelled [`type: enhancement`][enhancement].

Please avoid:

- Opening pull requests for issues marked as `blocked`.

All enhancement and bug issues are marked with a `level` label, it may help you know the
size/complexity of it.

## Building the project

Prerequisites:
- Go 1.19+

Run `make deps-install` to install the packages used by the project.

Run `make deps-upgrade` if you need to upgrade all of them, run `go help get` to see how to update
individual ones.

To build your changes into a binary run `make dist`, all three versions (Windows, Mac and Linux)
will be created under the `dist/` folder.

You can also just run `go run cmd/clockify-cli/main.go` to execute the source directly.

See the [project layout documentation][project layout] to know where to find and create specific
components.

## Submitting a pull request

Contributions to this project are [released][legal] to the public under the
[project's open source license][license]. By participating in this project you agree to abide by
its terms.

We generate manual pages from source on every release. You do not need to submit pull requests for
documentation specifically; manual pages for commands will automatically get updated after your
pull requests gets accepted.

### With [`gh`][gh]

1. Clone this repository
2. Make and commit your changes.
3. Submit a pull request: `gh pr create --web`
4. In its body link which issue it is related, if there is one

### Without `gh`

1. [Fork the repository][fork]
2. Make and commit your changes
3. [Open a pull request][open-pr]
4. In its body link which issue it is related, if there is one

## Resources

- [How to Contribute to Open Source][]
- [Using Pull Requests][]
- [GitHub Help][]

## Credits

This document is based on the [CONTRIBUTING.md from github/cli/cli][credit].

[fork]: https://github.com/lucassabreu/clockify-cli/fork
[open-pr]: https://github.com/lucassabreu/clockify-cli/compare
[credit]: https://github.com/cli/cli/blob/trunk/.github/CONTRIBUTING.md
[issues]: https://github.com/lucassabreu/clockify-cli/issues
[bugs]: https://github.com/lucassabreu/clockify-cli/issues?q=is%3Aopen+is%3Aissue+label%3A%22type%3A+bug%22
[enhancement]: https://github.com/lucassabreu/clockify-cli/issues?q=is%3Aissue+is%3Aopen+label%3A%22type%3A+enhancement%22
[project layout]: ./docs/project-layout.md
[gh]: https://github.com/cli/cli
[legal]: https://docs.github.com/en/free-pro-team@latest/github/site-policy/github-terms-of-service#6-contributions-under-repository-license
[license]: ./LICENSE
[How to Contribute to Open Source]: https://opensource.guide/how-to-contribute/
[Using Pull Requests]: https://docs.github.com/en/free-pro-team@latest/github/collaborating-with-issues-and-pull-requests/about-pull-requests
[GitHub Help]: https://docs.github.com/
