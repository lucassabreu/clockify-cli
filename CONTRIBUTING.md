## Contributing

Hi! Thanks for your interest in contributing to the Clockify CLI!

We accept pull requests for bug fixes and features where we've discussed the
approach in an issue and given the go-ahead on it. We'd also love to hear about
ideas for new features as issues.

Please do:

* Check existing issues to verify that the [bug][bug issues],
  [improvement][improvement issues] or [feature request][feature request issues]
  has not already been submitted.
* Open an issue if things aren't working as expected.
* Open an issue to propose a significant change.
* Open a pull request to fix a bug.
* Open a pull request to fix documentation about a command.
* Open a pull request for any issue labelled [`help wanted`][hw] or
  [`good first issue`][gfi].

Please avoid:

* Opening pull requests for issues marked `blocked`, `cancelled`, or
  `duplicate`.

## Building the project

Prerequisites:
- Go 1.17+

Build with:
* Unix-like systems: `make build-linux`
* Windows: `make build-windows`

Run the new binary as:
* Unix-like systems: `dist/linux/clockify-cli`
* Windows: `dist\windows\clockify-cli`

See [project layout documentation](docs/project-layout.md) for information on
where to find specific source files.

## Submitting a pull request

1. Create a new branch: `git checkout -b my-branch-name`
1. Make your change and test it locally
1. Submit the pull request linking it to the issue

Contributions to this project are [released][legal] to the public under the
[project's open source license][license].

We generate manual pages from source on every release. You do not need to
submit pull requests for documentation specifically; manual pages for commands
will automatically get updated after your pull requests gets accepted.

## Design guidelines

TODO: there is one?

## Resources

- [How to Contribute to Open Source][]
- [Using Pull Requests][]
- [Clockify API][]

## Credits

This document is based on the [CONTRIBUTING.md from github/cli/cli][credit].

[bug issues]: https://github.com/lucassabreu/clockify-cli/labels/type%3A%20bug
[feature request issues]: https://github.com/lucassabreu/clockify-cli/labels/type%3A%20new%20feature
[improvement issues]: https://github.com/lucassabreu/clockify-cli/labels/type%3A%20improvement
[hw]: https://github.com/lucassabreu/clockify-cli/labels/help%20wanted
[gfi]: https://github.com/lucassabreu/clockify-cli/labels/good%20first%20issue
[legal]: https://docs.github.com/en/free-pro-team@latest/github/site-policy/github-terms-of-service#6-contributions-under-repository-license
[license]: LICENSE
[code-of-conduct]: ./CODE-OF-CONDUCT.md
[How to Contribute to Open Source]: https://opensource.guide/how-to-contribute/
[Using Pull Requests]: https://docs.github.com/en/free-pro-team@latest/github/collaborating-with-issues-and-pull-requests/about-pull-requests
[Clockify API]: https://clockify.me/developers-api
[CLI Design System]: https://primer.style/cli/
[Google Docs Template]: https://docs.google.com/document/d/1JIRErIUuJ6fTgabiFYfCH3x91pyHuytbfa0QLnTfXKM/edit#heading=h.or54sa47ylpg
[credit]: https://github.com/cli/cli/blob/8f3b6749d7dd7ceafa1a15d211a5cb4d32422b22/.github/CONTRIBUTING.md
