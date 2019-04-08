workflow "Lint golang" {
  on = "push"
  resolves = ["sjkaliski/go-github-actions/lint@v0.2.0"]
}

action "sjkaliski/go-github-actions/lint@v0.2.0" {
  uses = "sjkaliski/go-github-actions/lint@v0.2.0"
  secrets = ["GITHUB_TOKEN"]
}

workflow "Release" {
  on = "push"
  resolves = ["goreleaser"]
}

workflow "Release Next" {
  on = "push"
  resolves = ["goreleaser-next"]
}

action "is-tag" {
  uses = "actions/bin/filter@master"
  args = "tag"
}

action "is-master" {
  uses = "actions/bin/filter@master"
  args = "branch master"
}

action "goreleaser" {
  uses = "docker://goreleaser/goreleaser"
  secrets = [
    "GITHUB_TOKEN",
    "DOCKER_USERNAME",
    "DOCKER_PASSWORD",
  ]
  args = "release"
  needs = ["is-tag"]
}

action "goreleaser-next" {
  uses = "docker://goreleaser/goreleaser"
  secrets = [
    "GITHUB_TOKEN",
    "DOCKER_USERNAME",
    "DOCKER_PASSWORD",
  ]
  args = "release --snapshot"
  needs = ["is-master"]
}
