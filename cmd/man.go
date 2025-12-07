package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/spf13/cobra/doc"
)

var manDir string

var manCmd = &cobra.Command{
    Use:   "man",
    Short: "Generate man pages for penguinguide",
    RunE: func(cmd *cobra.Command, args []string) error {
        if manDir == "" {
            manDir = "./man"
        }

        if err := os.MkdirAll(manDir, 0o755); err != nil {
            return err
        }

        header := &doc.GenManHeader{
            Title:   "PENGUINGUIDE",
            Section: "1",
        }

        fmt.Println("Writing man pages to", manDir)

        return doc.GenManTree(RootCmd, header, manDir)
    },
}

func init() {
    manCmd.Flags().StringVar(&manDir, "dir", "", "output directory for man pages")
    RootCmd.AddCommand(manCmd)
}

