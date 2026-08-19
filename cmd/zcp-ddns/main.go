package main

import (
	"fmt"
	"os"

	"github.com/zsoftly/zcp-ddns/internal/version"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "--version", "-v":
			fmt.Println(version.Version)
			return
		case "help", "--help", "-h":
			printUsage()
			return
		}
	}

	printUsage()
}

func printUsage() {
	fmt.Fprintf(os.Stdout, `zcp-ddns %s

Dynamic DNS service for ZSoftly DNS.

This repository is in initial development. See README.md and GitHub issues for
the implementation roadmap.
`, version.Version)
}
