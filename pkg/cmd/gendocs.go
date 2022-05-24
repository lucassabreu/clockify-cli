package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const gendocFrontmatterTemplate = `---
date: %s
title: "%s"
slug: %s
url: %s
weight: 40
---
`

// gendocsCmd represents the gendocs command
var gendocsCmd = &cobra.Command{
	Use:   "gendocs <output-dir>",
	Short: "Generate Markdown documentation for the clockify-cli.",
	Long: `Generate Markdown documentation for the clockify-cli.
This command is, mostly, used to create up-to-date documentation
of the command-line interface.

It creates one Markdown file per command with front matter suitable
for rendering in Hugo.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		docdir := "site/content/en/commands"
		if len(args) >= 1 {
			docdir = args[0]
		}

		if err := os.MkdirAll(docdir, os.ModePerm); err != nil {
			return err
		}

		now := time.Now().Format("2006-01-02")
		prepender := func(filename string) string {
			name := filepath.Base(filename)
			base := strings.TrimSuffix(name, path.Ext(name))
			url := "/en/commands/" + strings.ToLower(base) + "/"
			return fmt.Sprintf(gendocFrontmatterTemplate, now, strings.ReplaceAll(base, "_", " "), base, url)
		}

		linkHandler := func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))
			return "/en/commands/" + strings.ToLower(base) + "/"
		}

		fmt.Println("Generating Hugo command-line documentation in", docdir, "...")
		err := doc.GenMarkdownTreeCustom(cmd.Root(), docdir, prepender, linkHandler)
		if err != nil {
			return err
		}
		fmt.Println("Done.")

		return nil
	},
}
