package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"

    "penguinguide/internal/distro"
    "penguinguide/internal/ui"
)

var detectCmd = &cobra.Command{
    Use:   "detect",
    Short: "Detect the Linux distribution",
    Run: func(cmd *cobra.Command, args []string) {
        runDetect()
    },
}

func init() {
    RootCmd.AddCommand(detectCmd)
}

func runDetect() {
    d, err := distro.Detect()
    if err != nil {
        fmt.Fprintln(os.Stderr, ui.Error("Could not detect distribution"))
        fmt.Fprintln(os.Stderr, "  Error:", err)
        os.Exit(1)
    }

    fmt.Println(ui.Heading("Detected Linux distribution"))
    fmt.Printf("  %s %s\n", ui.Key("ID        :"), ui.Value(d.ID))
    fmt.Printf("  %s %v\n", ui.Key("ID_LIKE   :"), d.IDLike)
    fmt.Printf("  %s %s\n", ui.Key("NAME      :"), ui.Value(d.Name))
    fmt.Printf("  %s %s\n", ui.Key("PRETTY    :"), ui.Value(d.PrettyName))
    fmt.Printf("  %s %s\n", ui.Key("VERSION   :"), ui.Value(d.VersionID))
    fmt.Printf("  %s %s\n", ui.Key("FAMILY    :"), ui.Value(string(d.Family)))
}

