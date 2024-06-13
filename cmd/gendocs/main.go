package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/pkg/cmd"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
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

func main() {
	if err := execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func execute() error {
	docdir := "site/content/commands"
	if len(os.Args) > 1 {
		docdir = os.Args[1]
	}

	if err := os.MkdirAll(docdir, os.ModePerm); err != nil {
		return err
	}

	now := time.Now().Format("2006-01-02")
	prepender := func(filename string) string {
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		url := "/en/commands/" + strings.ToLower(base) + "/"
		return fmt.Sprintf(gendocFrontmatterTemplate,
			now, strings.ReplaceAll(base, "_", " "), base, url)
	}

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return "/en/commands/" + strings.ToLower(base) + "/"
	}

	cmd := cmd.NewCmdRoot(cmdutil.NewFactory(cmdutil.Version{}))

	fmt.Println("Generating Hugo command-line documentation in", docdir, "...")
	err := doc.GenMarkdownTreeCustom(cmd, docdir, prepender, linkHandler)
	if err != nil {
		return err
	}
	fmt.Println("Done.")

	return nil
}
