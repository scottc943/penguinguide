package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"

    "penguinguide/internal/distro"
    "penguinguide/internal/pkgmgr"
    "penguinguide/internal/ui"
)

var updateCmd = &cobra.Command{
    Use:   "update",
    Short: "Update installed packages",
    Run: func(cmd *cobra.Command, args []string) {
        runUpdate()
    },
}

func init() {
    RootCmd.AddCommand(updateCmd)
}

func runUpdate() {
    d, err := distro.Detect()
    if err != nil {
        fmt.Fprintln(os.Stderr, ui.Error("Could not detect distribution"))
        fmt.Fprintln(os.Stderr, "  Error:", err)
        os.Exit(1)
    }

    fmt.Println(ui.Heading("Update packages"))
    fmt.Printf("  %s %s\n", ui.Key("Distro family:"), ui.Value(string(d.Family)))

    mgr := pkgmgr.New(d)
    opts := pkgmgr.Options{
        DryRun:    dryRun,
        AssumeYes: assumeYes,
        Explain:   explain,
    }

    if err := mgr.UpdateAll(opts); err != nil {
        fmt.Fprintln(os.Stderr)
        fmt.Fprintln(os.Stderr, ui.Error("Update did not complete successfully"))
        fmt.Fprintln(os.Stderr, ui.Muted("The package manager output above has the detail"))
        os.Exit(1)
    }

    fmt.Println(ui.Success("Update finished"))
}

