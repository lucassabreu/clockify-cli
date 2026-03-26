package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func bumpVersion(v Version, bumpType string) Version {
	switch bumpType {
	case "major":
		return Version{Major: v.Major + 1, Minor: 0, Patch: 0}
	case "minor":
		return Version{Major: v.Major, Minor: v.Minor + 1, Patch: 0}
	case "patch":
		return Version{Major: v.Major, Minor: v.Minor, Patch: v.Patch + 1}
	default:
		return v
	}
}

func findLatestVersion(content string) (Version, error) {
	re := regexp.MustCompile(`## \[v(\d+)\.(\d+)\.(\d+)\]`)
	matches := re.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return Version{}, fmt.Errorf("no version found in CHANGELOG.md")
	}
	first := matches[0]
	major, _ := strconv.Atoi(first[1])
	minor, _ := strconv.Atoi(first[2])
	patch, _ := strconv.Atoi(first[3])
	return Version{Major: major, Minor: minor, Patch: patch}, nil
}

func findUnreleasedSection(content string) (string, int, int) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var lines []string
	inUnreleased := false
	startLine := -1
	endLine := -1

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

		if strings.HasPrefix(line, "## [Unreleased]") {
			inUnreleased = true
			startLine = len(lines) - 1
			continue
		}

		if inUnreleased && strings.HasPrefix(line, "## [") {
			endLine = len(lines) - 1
			break
		}
	}

	if startLine == -1 {
		return "", -1, -1
	}

	if endLine == -1 {
		endLine = len(lines)
	}

	section := strings.Join(lines[startLine:endLine], "\n")
	return section, startLine, endLine
}

func extractUnreleasedContent(unreleasedSection string) string {
	lines := strings.Split(unreleasedSection, "\n")
	if len(lines) <= 1 {
		return ""
	}

	var content []string
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "## [") {
			break
		}
		content = append(content, line)
	}

	return strings.TrimSpace(strings.Join(content, "\n"))
}

func updateUnreleasedLink(content string, newVersion string) string {
	oldUnreleasedLink := regexp.MustCompile(`\[Unreleased\]:.*`).FindString(content)
	if oldUnreleasedLink == "" {
		return content
	}

	newLink := fmt.Sprintf("[Unreleased]: https://github.com/lucassabreu/clockify-cli/compare/%s...HEAD", newVersion) +
		"\n" +
		fmt.Sprintf("[%s]: https://github.com/lucassabreu/clockify-cli/releases/tag/%s", newVersion, newVersion)
	return strings.Replace(content, oldUnreleasedLink, newLink, 1)
}

func runGitCommand(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "GIT_EDITOR=true", "GIT_SEQUENCE_EDITOR=true")
	return cmd.Run()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: release major|minor|patch")
		os.Exit(1)
	}

	bumpType := strings.ToLower(os.Args[1])
	if bumpType != "major" && bumpType != "minor" && bumpType != "patch" {
		fmt.Fprintln(os.Stderr, "Error: argument must be major, minor, or patch")
		os.Exit(1)
	}

	changelogPath := "CHANGELOG.md"
	if _, err := os.Stat(changelogPath); os.IsNotExist(err) {
		changelogPath = "../../CHANGELOG.md"
	}

	contentBytes, err := os.ReadFile(changelogPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", changelogPath, err)
		os.Exit(1)
	}
	content := string(contentBytes)

	unreleasedSection, startLine, _ := findUnreleasedSection(content)
	if unreleasedSection == "" || startLine == -1 {
		fmt.Fprintln(os.Stderr, "Error: no [Unreleased] section found in CHANGELOG.md")
		os.Exit(1)
	}

	unreleasedContent := extractUnreleasedContent(unreleasedSection)
	if unreleasedContent == "" {
		fmt.Fprintln(os.Stderr, "Error: [Unreleased] section is empty, nothing to release")
		os.Exit(1)
	}

	latestVersion, err := findLatestVersion(content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding latest version: %v\n", err)
		os.Exit(1)
	}

	newVersion := bumpVersion(latestVersion, bumpType)
	date := time.Now().Format("2006-01-02")

	fmt.Printf("Releasing %s (bumping %s from %s)\n", newVersion, bumpType, latestVersion)

	newVersionSection := fmt.Sprintf("## [%s] - %s\n\n%s\n", newVersion, date, unreleasedContent)

	lines := strings.Split(content, "\n")

	var unreleasedStartIdx, unreleasedEndIdx int
	for i, line := range lines {
		if strings.HasPrefix(line, "## [Unreleased]") {
			unreleasedStartIdx = i
		}
		if unreleasedStartIdx > 0 && strings.HasPrefix(line, "## [v") {
			unreleasedEndIdx = i
			break
		}
	}

	var newLines []string
	newLines = append(newLines, lines[:unreleasedStartIdx+1]...)
	newLines = append(newLines, "")
	newLines = append(newLines, newVersionSection)
	newLines = append(newLines, lines[unreleasedEndIdx:]...)

	newContent := strings.Join(newLines, "\n")

	newContent = updateUnreleasedLink(newContent, newVersion.String())

	if err := os.WriteFile(changelogPath, []byte(newContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", changelogPath, err)
		os.Exit(1)
	}

	fmt.Printf("Updated %s\n", changelogPath)

	if err := runGitCommand("git", "add", changelogPath); err != nil {
		fmt.Fprintln(os.Stderr, "Error: git add failed")
		os.Exit(1)
	}

	commitMsg := fmt.Sprintf("release: %s", newVersion)
	if err := runGitCommand("git", "commit", "-m", commitMsg); err != nil {
		fmt.Fprintln(os.Stderr, "Error: git commit failed")
		os.Exit(1)
	}
	fmt.Printf("Created commit: %s\n", commitMsg)

	if err := runGitCommand("git", "tag", "-a", newVersion.String(), "-m", fmt.Sprintf("Release %s", newVersion)); err != nil {
		fmt.Fprintln(os.Stderr, "Error: git tag failed")
		os.Exit(1)
	}
	fmt.Printf("Created tag: %s\n", newVersion)

	fmt.Println("Release completed successfully!")
}
