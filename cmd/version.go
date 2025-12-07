package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "penguinguide/internal/ui"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Show penguinguide version information",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println(ui.Heading("penguinguide version"))
        fmt.Printf("Version: %s\n", ui.Value(buildVersion))
        fmt.Printf("Commit : %s\n", ui.Value(buildCommit))
        fmt.Printf("Built  : %s\n", ui.Value(buildDate))
    },
}

func init() {
    RootCmd.AddCommand(versionCmd)
}

