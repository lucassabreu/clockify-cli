workflow "Lint golang" {
  on = "push"
  resolves = ["sjkaliski/go-github-actions/lint@v0.2.0"]
}

action "sjkaliski/go-github-actions/lint@v0.2.0" {
  uses = "sjkaliski/go-github-actions/lint@v0.2.0"
  secrets = ["GITHUB_TOKEN"]
}
