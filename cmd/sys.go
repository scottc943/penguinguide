package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"

    "penguinguide/internal/sysinfo"
    "penguinguide/internal/ui"
)

var sysCmd = &cobra.Command{
    Use:   "sys",
    Short: "Show system information",
    Long: `Show a quick summary of this system.

By default this prints an overview. Subcommands show
details such as network and WiFi.`,
    Run: func(cmd *cobra.Command, args []string) {
        runSysSummary()
    },
}

func init() {
    RootCmd.AddCommand(sysCmd)
}

func runSysSummary() {
    summary, err := sysinfo.GetSystemSummary()
    if err != nil {
        fmt.Fprintln(os.Stderr, ui.Error("Could not read system information"))
        fmt.Fprintln(os.Stderr, "  Error:", err)
        os.Exit(1)
    }

    fmt.Println(ui.Heading("System summary"))
    fmt.Printf("  %s %s\n", ui.Key("Hostname     :"), ui.Value(summary.Hostname))
    fmt.Printf("  %s %s\n", ui.Key("Distribution :"), ui.Value(summary.DistroName))
    fmt.Printf("  %s %s\n", ui.Key("Kernel       :"), ui.Value(summary.Kernel))
    fmt.Printf("  %s %s\n", ui.Key("Uptime       :"), ui.Value(summary.Uptime))
    fmt.Printf("  %s %s\n", ui.Key("Load average :"), ui.Value(summary.LoadAverage))
    fmt.Printf("  %s %s\n", ui.Key("Memory usage :"), ui.Value(summary.MemoryPretty))
}

