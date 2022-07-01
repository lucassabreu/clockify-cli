![Clockify CLI](https://repository-images.githubusercontent.com/173476481/3445a278-9bb9-49e9-8c99-d10c76574489)
============

A simple cli to manage your time entries and projects on Clockify from terminal

[![Release](https://img.shields.io/github/release/lucassabreu/clockify-cli.svg?classes=badges)](https://github.com/lucassabreu/clockify-cli/releases/latest)
[![GitHub all releases](https://img.shields.io/github/downloads/lucassabreu/clockify-cli/total)](https://github.com/lucassabreu/clockify-cli/releases)
[![clockify-cli](https://snapcraft.io//clockify-cli/badge.svg?classes=badges)](https://snapcraft.io/clockify-cli)
[![Build Status](https://github.com/lucassabreu/clockify-cli/actions/workflows/release.yml/badge.svg?classes=badges)](.github/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/lucassabreu/clockify-cli?classes=badges)](https://goreportcard.com/report/github.com/lucassabreu/clockify-cli)
[![Netlify Status](https://api.netlify.com/api/v1/badges/8667b9f6-4ca2-4ee4-865e-20b5848e7059/deploy-status?classes=badges)](https://app.netlify.com/sites/clockify-cli/deploys)
[![DeepSource](https://deepsource.io/gh/lucassabreu/clockify-cli.svg/?classes=badges&label=active+issues&show_trend=true&token=hkvNbnaRCE4DhtW6vDYpFWSR)](https://deepsource.io/gh/lucassabreu/clockify-cli/?ref=repository-badge)

Documentation
-------------

See the [project site](https://clockify-cli.netlify.app/) for the how to setup and use this CLI.

See more information about the sub-commands at: https://clockify-cli.netlify.app/en/commands/clockify-cli/

Contributing
------------

On how to help improve the tool, suggest new features or report bugs, please take a look at the
[CONTRIBUTING.md](CONTRIBUTING.md).

Features
--------

* [x] List time entries from a day
  * [x] List in progress entry
* [x] Report time entries using a date range
  * [x] Inform date range as parameters
  * [x] "auto filter" for last month
  * [x] "auto filter" for this month
* [x] Start a new time entry
  * [x] Cloning last time entry
  * [x] Ask input interactively
* [x] Stop the last entry
* [x] List workspace projects
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

