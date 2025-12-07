package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"

    "penguinguide/internal/distro"
    "penguinguide/internal/pkgmgr"
    "penguinguide/internal/ui"
)

var infoCmd = &cobra.Command{
    Use:   "info [package]",
    Short: "Show detailed information for a package",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        runInfo(args[0])
    },
}

func init() {
    RootCmd.AddCommand(infoCmd)
}

func runInfo(name string) {
    d, err := distro.Detect()
    if err != nil {
        fmt.Fprintln(os.Stderr, ui.Error("Could not detect distribution"))
        fmt.Fprintln(os.Stderr, "  Error:", err)
        os.Exit(1)
    }

    fmt.Println(ui.Heading("Package information"))
    fmt.Printf("  %s %s\n", ui.Key("Distro family:"), ui.Value(string(d.Family)))
    fmt.Printf("  %s %s\n", ui.Key("Package      :"), ui.Value(name))
    fmt.Println()

    mgr := pkgmgr.New(d)

    opts := pkgmgr.Options{
        DryRun:    dryRun,
        AssumeYes: assumeYes,
        Explain:   explain,
    }

    if err := mgr.Info(name, opts); err != nil {
        fmt.Fprintln(os.Stderr)
        fmt.Fprintln(os.Stderr, ui.Error("Package info request did not complete successfully"))
        fmt.Fprintln(os.Stderr, ui.Muted("The package manager output above has the detail"))
        os.Exit(1)
    }
}

