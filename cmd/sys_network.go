package cmd

import (
    "bufio"
    "fmt"
    "net/http"
    "os"
    "strings"

    "github.com/spf13/cobra"

    "penguinguide/internal/sysinfo"
    "penguinguide/internal/ui"
)

var sysNetCmd = &cobra.Command{
    Use:   "network",
    Short: "Show network configuration with explanations",
    Run: func(cmd *cobra.Command, args []string) {
        runSysNetwork()
    },
}

func init() {
    sysCmd.AddCommand(sysNetCmd)
}

func runSysNetwork() {
    fmt.Println(ui.Heading("Network overview"))
    fmt.Println()

    iface, gw := sysinfo.GetDefaultRoute()
    if gw != "" {
        fmt.Printf("  %s %s via %s\n",
            ui.Key("Default gateway:"), ui.Value(gw), ui.Value(iface))
    } else {
        fmt.Printf("  %s %s\n", ui.Key("Default gateway:"), ui.Warning("unknown"))
    }

    dnsServers := sysinfo.GetDNSServers()
    fmt.Printf("  %s\n", ui.Key("DNS servers:"))
    if len(dnsServers) == 0 {
        fmt.Println("    " + ui.Warning("none found"))
    } else {
        for _, s := range dnsServers {
            fmt.Println("    " + ui.Value(s))
        }
    }

    fmt.Println()
    fmt.Print(ui.Info("Check public IP address") + " [y/N]: ")

    var ans string
    fmt.Scanln(&ans)
    ans = strings.ToLower(strings.TrimSpace(ans))

    if ans == "y" || ans == "yes" {
        ip := fetchPublicIP()
        fmt.Printf("  %s %s\n", ui.Key("Public IP     :"), ui.Value(ip))
        fmt.Println()
    }

    fmt.Println(ui.Heading("Interfaces"))

    infos, err := sysinfo.GetInterfaceInfo()
    if err != nil {
        fmt.Fprintln(os.Stderr, ui.Error("Could not read interfaces"))
        fmt.Fprintln(os.Stderr, "  Error:", err)
        os.Exit(1)
    }

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

func fetchPublicIP() string {
    resp, err := http.Get("https://api.ipify.org")
    if err != nil {
        return "unknown"
    }
    defer resp.Body.Close()

    scanner := bufio.NewScanner(resp.Body)
    if scanner.Scan() {
        return scanner.Text()
    }
    return "unknown"
}

