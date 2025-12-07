package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"

    "penguinguide/internal/ui"
)

var (
    dryRun    bool
    assumeYes bool
    explain   bool
)

var RootCmd = &cobra.Command{
    Use:   "penguinguide",
    Short: "Friendly helper for Linux newcomers",
    Long: `penguinguide explains what your system is doing
and shows the native commands behind each action.`,
}

func Execute() {
    if err := RootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, ui.Error("Error: "), err)
        os.Exit(1)
    }
}

func init() {
    RootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", true, "show commands before running them and ask for confirmation")
    RootCmd.PersistentFlags().BoolVarP(&assumeYes, "yes", "y", false, "assume yes for package operations")
    RootCmd.PersistentFlags().BoolVar(&explain, "explain", false, "explain what penguinguide is doing and show native commands")
}

