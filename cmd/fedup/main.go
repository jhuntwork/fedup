package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/jhuntwork/fedup"
	"github.com/spf13/cobra"
)

const (
	name        = "fedup"
	directory   = " <directory>"
	description = "\nfedup is a simple File dEDUPlicator.\n\n" +
		"It walks a given directory, collects hash sums of all files,\n" +
		"and turns all duplicates into hard links."
)

var (
	dryrun bool
	quiet  bool
)

func dedup(cmd *cobra.Command, args []string) {
	dirname := args[0]
	var out io.Writer
	if quiet {
		out = &bytes.Buffer{}
	} else {
		out = os.Stdout
	}
	if dryrun {
		fmt.Fprintf(out, "Dry run mode: no changes being made.\n\n")
	}
	count, err := fedup.Dedup(dirname, dryrun, out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(out, "Found %d duplicates\n", count)
}

func main() {
	rootCmd := &cobra.Command{
		Use:   name + directory,
		Args:  cobra.ExactArgs(1),
		Short: description,
		Run:   dedup,
	}
	rootCmd.PersistentFlags().BoolVarP(&dryrun, "dryrun", "d", false, "Do not make any changes.")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Do not print any output.")
	if err := rootCmd.Execute(); err != nil {
		return
	}
}
