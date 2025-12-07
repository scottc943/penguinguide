package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"

    "penguinguide/internal/distro"
    "penguinguide/internal/pkgmgr"
    "penguinguide/internal/ui"
)

var installCmd = &cobra.Command{
    Use:   "install [packages...]",
    Short: "Install packages",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        runInstall(args)
    },
}

func init() {
    RootCmd.AddCommand(installCmd)
}

func runInstall(pkgs []string) {
    d, err := distro.Detect()
    if err != nil {
        fmt.Fprintln(os.Stderr, ui.Error("Could not detect distribution"))
        fmt.Fprintln(os.Stderr, "  Error:", err)
        os.Exit(1)
    }

    fmt.Println(ui.Heading("Install packages"))
    fmt.Printf("  %s %s\n", ui.Key("Distro family:"), ui.Value(string(d.Family)))
    fmt.Printf("  %s %v\n", ui.Key("Packages     :"), pkgs)
    fmt.Println()

    mgr := pkgmgr.New(d)
    opts := pkgmgr.Options{
        DryRun:    dryRun,
        AssumeYes: assumeYes,
        Explain:   explain,
    }

    if err := mgr.Install(pkgs, opts); err != nil {
        fmt.Fprintln(os.Stderr)
        fmt.Fprintln(os.Stderr, ui.Error("Package install did not complete successfully"))
        fmt.Fprintln(os.Stderr, ui.Muted("Most of the time this means the package manager reported an error"))
        os.Exit(1)
    }

    fmt.Println(ui.Success("Install finished"))
}

