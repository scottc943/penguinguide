package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"

    "penguinguide/internal/sysinfo"
    "penguinguide/internal/ui"
)

var sysIPCmd = &cobra.Command{
    Use:   "ip",
    Short: "Show network interfaces and IP addresses",
    Run: func(cmd *cobra.Command, args []string) {
        runSysIP()
    },
}

func init() {
    sysCmd.AddCommand(sysIPCmd)
}

func runSysIP() {
    infos, err := sysinfo.GetInterfaceInfo()
    if err != nil {
        fmt.Fprintln(os.Stderr, ui.Error("Could not read network interface information"))
        fmt.Fprintln(os.Stderr, "  Error:", err)
        os.Exit(1)
    }

    if len(infos) == 0 {
        fmt.Println(ui.Warning("No network interfaces with addresses found"))
        return
    }

    fmt.Println(ui.Heading("Network interfaces"))

    for _, iface := range infos {
        state := ui.Value("down")
        if iface.IsUp {
            state = ui.Success("up")
        }
        loop := ""
        if iface.IsLoopback {
            loop = " " + ui.Muted("(loopback)")
        }

        fmt.Printf("  %s (%s)%s\n", ui.Value(iface.Name), state, loop)
        for _, addr := range iface.Addresses {
            fmt.Println("    " + ui.Info(addr))
        }
    }
}

