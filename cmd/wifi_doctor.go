package cmd

import (
    "fmt"

    "github.com/spf13/cobra"

    "penguinguide/internal/ui"
)

var wifiDoctorCmd = &cobra.Command{
    Use:   "wifi-doctor",
    Short: "Run WiFi checks and a quick speedtest",
    Long: `wifi-doctor runs WiFi diagnostics and a speedtest in sequence.

It first shows WiFi connection details, signal, band, channel,
security, and latency. Then it runs a smaller download test
so you can see practical performance.`,
    Run: func(cmd *cobra.Command, args []string) {
        runWifiDoctor()
    },
}

func init() {
    RootCmd.AddCommand(wifiDoctorCmd)
}

func runWifiDoctor() {
    fmt.Println(ui.Heading("WiFi doctor"))
    fmt.Println()
    fmt.Println(ui.Info("Step 1: WiFi connection, signal, and latency"))
    fmt.Println()

    runSysWifiNonInteractive()

    fmt.Println()
    fmt.Println(ui.Info("Step 2: Quick speed test"))
    fmt.Println()

    runSpeedtestQuickNonInteractive()

    fmt.Println()
    fmt.Println(ui.Success("WiFi doctor finished"))
    fmt.Println(ui.Muted("Use the suggestions, latency, and speed results above to choose next steps."))
}

