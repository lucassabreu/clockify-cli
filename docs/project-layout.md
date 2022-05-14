# Clockify CLI Project Layout

The project is organized in the following folders and important files:

- [`cmd/`](../cmd) - `main` packages to build executable binaries.
- [`docs/`](.) - documentation of the project for maintainers and contributors.
- [`scripts/`](../scripts) - build and release scripts.
- [`api/`](../api) - golang implementation of the Clockify API.
- [`pkg/`](../pkg) - other packages that support the `api` or commands.
- [`internal/`](../internal) - Go packages that are highly specific to this project
- [`go.mod`](../go.mod) - external Go dependencies for this project.
- [`Makefile`](../Makefile) - most of setup and maintenance actions for this project.

## Command line organization

All CLI commands will be under [`pkg/cmd/`](../pkg/cmd) and the file naming convention is this:

```
pkg/cmd/<command>/<subcommand>/<subcommand>.go
```

Following the same structure as the command path, so `clockify-cli client add` is at
`pkg/cmd/client/add.go`, all command packages will have a function named `NewCmd<subcommand>` that
will receive a `*cmdutil.Factory` type and return a `*cobra.Command`.

Specific logic for that command must be kept at the same package as it, and all subcommands must be
registered on its parent package. So all subcommands of `client` will registered on the function
`client.NewCmdClient()`.

Output formatters must stay under the package [`pkg/output/`](../pkg/output) using the following
file convention:

```
pkg/output/<entity>/<format>.go
```

Shared functionality for printing entities must be at the package
[`pkg/outpututil/`](../pkg/outpututil).

### Steps do create a new command

Say you will create a new command `delete` under `client`.

1. Create the package `pkg/cmd/client/delete/`
2. Create a function called `NewCmdDelete` on a file `delete.go`
    1. This function must receive a [`*cmdutil.Factory`][] struct and
       return a [`*cobra.Command`][] fully set.
3. Edit the entity root command at `pkg/cmd/client/client.go` to register it as a subcommand using
   the factory function previously created. If the entity root does not exist yet, then:
    1. Create the file, and in it a function `NewCmdClient` that should receive `*cmdutil.Factory`
       and return a [`*cobra.Command`][] with all its subcommands.
4. If is the first command of a entity:
    1. Create a package called `pkg/output/client`
    2. Implement output the five basic output formats `table` (default), `json`, `quiet` (only the
       ID), `template` ([Go template](https://pkg.go.dev/text/template)) and `csv`. Each one on a
       file by itself.

## Credits

This document is based on the [project-layout.md from github/cli/cli][credit].

[credit]: https://github.com/cli/cli/blob/trunk/docs/project-layout.md
[`*cobra.Command`]: https://pkg.go.dev/github.com/spf13/cobra#Command
[`*cmdutil.Factory`]: ../pkg/cmdutil/factory.go
