package cmd

import (
    "fmt"
    "os"
    "strings"

    "github.com/spf13/cobra"

    "penguinguide/internal/distro"
    "penguinguide/internal/pkgmgr"
    "penguinguide/internal/ui"
)

var searchCmd = &cobra.Command{
    Use:   "search [query...]",
    Short: "Search for packages by name or description",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        runSearch(args)
    },
}

func init() {
    RootCmd.AddCommand(searchCmd)
}

func runSearch(args []string) {
    query := strings.Join(args, " ")

    d, err := distro.Detect()
    if err != nil {
        fmt.Fprintln(os.Stderr, ui.Error("Could not detect distribution"))
        fmt.Fprintln(os.Stderr, "  Error:", err)
        os.Exit(1)
    }

    fmt.Println(ui.Heading("Package search"))
    fmt.Printf("  %s %s\n", ui.Key("Distro family:"), ui.Value(string(d.Family)))
    fmt.Printf("  %s %q\n", ui.Key("Query        :"), query)
    fmt.Println()

    mgr := pkgmgr.New(d)

    opts := pkgmgr.Options{
        DryRun:    dryRun,
        AssumeYes: assumeYes,
        Explain:   explain,
    }

    if err := mgr.Search(query, opts); err != nil {
        fmt.Fprintln(os.Stderr)
        fmt.Fprintln(os.Stderr, ui.Error("Package search did not complete successfully"))
        fmt.Fprintln(os.Stderr, ui.Muted("The package manager output above has the detail"))
        os.Exit(1)
    }
}

