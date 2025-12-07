package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"

    "penguinguide/internal/distro"
    "penguinguide/internal/pkgmgr"
    "penguinguide/internal/ui"
)

var removeCmd = &cobra.Command{
    Use:   "remove [packages...]",
    Short: "Remove packages",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        runRemove(args)
    },
}

func init() {
    RootCmd.AddCommand(removeCmd)
}

func runRemove(pkgs []string) {
    d, err := distro.Detect()
    if err != nil {
        fmt.Fprintln(os.Stderr, ui.Error("Could not detect distribution"))
        fmt.Fprintln(os.Stderr, "  Error:", err)
        os.Exit(1)
    }

    fmt.Println(ui.Heading("Remove packages"))
    fmt.Printf("  %s %s\n", ui.Key("Distro family:"), ui.Value(string(d.Family)))
    fmt.Printf("  %s %v\n", ui.Key("Packages     :"), pkgs)
    fmt.Println()

    mgr := pkgmgr.New(d)
    opts := pkgmgr.Options{
        DryRun:    dryRun,
        AssumeYes: assumeYes,
        Explain:   explain,
    }

    if err := mgr.Remove(pkgs, opts); err != nil {
        fmt.Fprintln(os.Stderr)
        fmt.Fprintln(os.Stderr, ui.Error("Package removal did not complete successfully"))
        fmt.Fprintln(os.Stderr, ui.Muted("The package manager output above has the detail"))
        os.Exit(1)
    }

    fmt.Println(ui.Success("Removal finished"))
}

