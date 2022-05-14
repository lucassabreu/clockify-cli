# Getting Started

Clockify CLI is a command line tool to help manage time entries from [Clockify][clockify] and
related resources.

- [Available commands][commands]
- [Usage examples][usage]
- [Installation][install]

## Configuration

- Generate a API key by visiting your user [settings on Clockify.me][settings], in the "API"
  section generate with you don't have one and copy your API key.
- Run the command `clockify-cli config --init`, it will ask for the API key you copied.
  - The CLI will ask for your preferences and default settings.
- After this run the command `clockify-cli in` to start a time entry.
- (optional) To add auto completion follow the instructions [here][auto-complete]

## How to Contact

- Questions about how to use the CLI?
- Wanna provide feedback on some feature?
- Report a bug or ask for a feature?

All these can be done opening a [issue on Github][issues].

## Contributing

Wants to help improve the CLI or the project, check out our [contributing page][contributing]

#### Disclaimer

The maintainers of this CLI are just users of Clockify and have no inside view from it, all actions
performed by it are possible using the [API][api] provided by Clockify.

[clockify]: https://clockify.me/
[api]: https://clockify.me/developers-api
[install]: https://github.com/lucassabreu/clockify-cli#how-to-install-
[usage]: /en/usage-examples/
[commands]: /en/commands/clockify-cli/
[settings]: https://app.clockify.me/user/settings
[auto-complete]: /en/commands/clockify-cli_completion/#synopsis
[issues]: https://github.com/lucassabreu/clockify-cli/issues
[contributing]: https://github.com/lucassabreu/clockify-cli/blob/main/CONTRIBUTING.md
