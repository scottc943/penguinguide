package cmd

import (
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"

    "github.com/spf13/cobra"

    "penguinguide/internal/sysinfo"
    "penguinguide/internal/ui"
)

var sysWifiCmd = &cobra.Command{
    Use:   "wifi",
    Short: "Show wireless connection details and helpful hints",
    Run: func(cmd *cobra.Command, args []string) {
        wifiCheck(true)
    },
}

func init() {
    sysCmd.AddCommand(sysWifiCmd)
}

type wifiStatus struct {
    Device        string
    SSID          string
    SignalPercent int
    QualityText   string
    Band          string
    FrequencyMHz  int
    FrequencyRaw  string
    Channel       int
    RateRaw       string
    SecurityRaw   string
}

// shared helper
func wifiCheck(interactive bool) {
    status, ok := getWifiStatus()
    if !ok {
        fmt.Println(ui.Error("WiFi info could not be determined"))
        fmt.Println(ui.Muted("You may need NetworkManager or wireless tools installed"))
        return
    }

    printWifiInfo(status)

    if interactive {
        fmt.Println()
        fmt.Print(ui.Info("Run a quick latency and packet loss test to 8.8.8.8") + " [y/N]: ")
        var ans string
        fmt.Fscan(os.Stdin, &ans)
        ans = strings.ToLower(strings.TrimSpace(ans))
        if ans == "y" || ans == "yes" {
            fmt.Println()
            avgMs, lossPct, err := runLatencyTest()
            if err != nil {
                fmt.Println(ui.Error("Latency test failed:"), err)
            } else {
                printLatencyInfo(avgMs, lossPct)
                fmt.Println()
                printWifiSuggestions(status, avgMs, lossPct)
                return
            }
        }
        fmt.Println()
        printWifiSuggestions(status, 0, 0)
        return
    }

    // non interactive for wifi doctor
    fmt.Println()
    avgMs, lossPct, err := runLatencyTest()
    if err != nil {
        fmt.Println(ui.Error("Latency test failed:"), err)
        fmt.Println()
        printWifiSuggestions(status, 0, 0)
        return
    }
    printLatencyInfo(avgMs, lossPct)
    fmt.Println()
    printWifiSuggestions(status, avgMs, lossPct)
}

func runSysWifiNonInteractive() {
    wifiCheck(false)
}

func getWifiStatus() (wifiStatus, bool) {
    if _, err := exec.LookPath("nmcli"); err == nil {
        if s, ok := wifiFromNmcli(); ok {
            return s, true
        }
    }
    if s, ok := wifiFromIwconfig(); ok {
        return s, true
    }
    return wifiStatus{}, false
}

func wifiFromNmcli() (wifiStatus, bool) {
    cmdStr := "nmcli -t -f ACTIVE,DEVICE,SSID,SIGNAL,FREQ,RATE,SECURITY dev wifi | grep '^yes:' 2>/dev/null"
    out, err := exec.Command("sh", "-c", cmdStr).Output()
    if err != nil || len(out) == 0 {
        return wifiStatus{}, false
    }

    line := strings.Split(strings.TrimSpace(string(out)), "\n")[0]
    fields := strings.Split(line, ":")
    if len(fields) < 7 {
        return wifiStatus{}, false
    }

    device := fields[1]
    ssid := strings.ReplaceAll(fields[2], `\:`, ":")
    signalStr := fields[3]
    freqStr := fields[4]
    rate := fields[5]
    security := fields[6]

    signalPercent, _ := strconv.Atoi(strings.TrimSpace(signalStr))
    band, freqMHz := sysinfo.BandFromFreq(freqStr)
    channel := sysinfo.ChannelFromFreq(freqMHz)

    return wifiStatus{
        Device:        device,
        SSID:          ssid,
        SignalPercent: signalPercent,
        Band:          band,
        FrequencyMHz:  freqMHz,
        FrequencyRaw:  freqStr,
        Channel:       channel,
        RateRaw:       rate,
        SecurityRaw:   security,
    }, true
}

func wifiFromIwconfig() (wifiStatus, bool) {
    out, err := exec.Command("sh", "-c", "iwconfig 2>/dev/null").Output()
    if err != nil || len(out) == 0 {
        return wifiStatus{}, false
    }

    lines := strings.Split(string(out), "\n")
    var ifaceLine string
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }
        if strings.Contains(line, "ESSID:") && !strings.Contains(line, "ESSID:off/any") {
            ifaceLine = line
            break
        }
    }

    if ifaceLine == "" {
        return wifiStatus{}, false
    }

    fields := strings.Fields(ifaceLine)
    if len(fields) == 0 {
        return wifiStatus{}, false
    }
    iface := fields[0]

    ssid := sysinfo.ExtractBetween(ifaceLine, `ESSID:"`, `"`)
    quality := sysinfo.ExtractAfter(ifaceLine, "Link Quality=")
    signal := sysinfo.ExtractAfter(ifaceLine, "Signal level=")

    freqStr := ""
    extraOut, err := exec.Command("sh", "-c", "iwconfig "+iface+" 2>/dev/null").Output()
    if err == nil {
        for _, l := range strings.Split(string(extraOut), "\n") {
            if strings.Contains(l, "Frequency:") {
                f := sysinfo.ExtractAfter(l, "Frequency:")
                freqStr = strings.TrimSpace(f)
                break
            }
        }
    }

    band, freqMHz := sysinfo.BandFromFreq(freqStr)
    signalPercent := sysinfo.QualityToPercent(quality)

    s := wifiStatus{
        Device:        iface,
        SSID:          ssid,
        SignalPercent: signalPercent,
        QualityText:   quality + " " + signal,
        Band:          band,
        FrequencyMHz:  freqMHz,
        FrequencyRaw:  freqStr,
        Channel:       sysinfo.ChannelFromFreq(freqMHz),
        RateRaw:       "",
        SecurityRaw:   "",
    }

    return s, true
}

func printWifiInfo(s wifiStatus) {
    fmt.Println(ui.Heading("WiFi connection"))

    fmt.Printf("  %s %s\n", ui.Key("Device   :"), ui.Value(s.Device))
    fmt.Printf("  %s %s\n", ui.Key("SSID     :"), ui.Value(safeValue(s.SSID)))

    if s.SignalPercent > 0 {
        desc := describeSignal(s.SignalPercent)
        col := colorForSignal(s.SignalPercent)
        fmt.Printf("  %s %s%d percent%s %s\n",
            ui.Key("Signal   :"), col, s.SignalPercent, ui.Reset, desc)
    } else if s.QualityText != "" {
        fmt.Printf("  %s %s\n", ui.Key("Signal   :"), s.QualityText)
    }

    if s.Band != "" {
        fmt.Printf("  %s %s\n", ui.Key("Band     :"), s.Band)
    }
    if s.Channel > 0 {
        hintText := sysinfo.ChannelHintText(s.Band, s.Channel)
        if hintText != "" {
            fmt.Printf("  %s %d %s\n", ui.Key("Channel  :"), s.Channel, ui.Muted(hintText))
        } else {
            fmt.Printf("  %s %d\n", ui.Key("Channel  :"), s.Channel)
        }
    }
    if s.FrequencyRaw != "" {
        fmt.Printf("  %s %s\n", ui.Key("Frequency:"), s.FrequencyRaw)
    }

    if s.RateRaw != "" {
        rateDesc := describeRate(s.RateRaw)
        fmt.Printf("  %s %s %s\n", ui.Key("Link speed:"), s.RateRaw, rateDesc)
    }

    if s.SecurityRaw != "" {
        secReadable := describeSecurity(s.SecurityRaw)
        fmt.Printf("  %s %s\n", ui.Key("Security :"), secReadable)
    }
}

func safeValue(v string) string {
    if strings.TrimSpace(v) == "" {
        return "(unknown)"
    }
    return v
}

func describeSignal(value int) string {
    switch {
    case value >= 80:
        return "(Great connection)"
    case value >= 60:
        return "(Good)"
    case value >= 40:
        return "(Fair)"
    case value >= 20:
        return "(Weak, slower speeds likely)"
    case value > 0:
        return "(Very weak, may disconnect)"
    default:
        return ""
    }
}

func colorForSignal(value int) string {
    switch {
    case value >= 70:
        return ui.Green
    case value >= 40:
        return ui.Yellow
    default:
        return ui.Red
    }
}

func describeSecurity(sec string) string {
    if sec == "--" || sec == "" {
        return ui.Error("Open (no password, not secure)")
    }

    desc := sec
    if strings.Contains(sec, "WPA3") {
        desc += " (modern and strong)"
    } else if strings.Contains(sec, "WPA2") {
        desc += " (good security for most networks)"
    } else if strings.Contains(sec, "WEP") {
        desc += " (very weak, should be avoided)"
    } else {
        desc += " (unknown style)"
    }
    return desc
}

func describeRate(rate string) string {
    if strings.Contains(rate, "Mbit") || strings.Contains(rate, "Mbps") {
        return ui.Muted("(WiFi link speed, not actual internet speed)")
    }
    return ""
}

func runLatencyTest() (avgMs float64, lossPct float64, err error) {
    cmd := exec.Command("ping", "-c", "4", "-w", "5", "8.8.8.8")
    out, err := cmd.CombinedOutput()
    if err != nil && len(out) == 0 {
        return 0, 0, err
    }

    text := string(out)
    lines := strings.Split(text, "\n")

    for _, line := range lines {
        line = strings.TrimSpace(line)
        if strings.Contains(line, "packet loss") {
            parts := strings.Split(line, ",")
            for _, p := range parts {
                p = strings.TrimSpace(p)
                if strings.Contains(p, "% packet loss") {
                    v := strings.Fields(p)[0]
                    v = strings.TrimSuffix(v, "%")
                    lossPct, _ = strconv.ParseFloat(v, 64)
                }
            }
        }
        if strings.Contains(line, "min/avg/max") {
            parts := strings.Split(line, "=")
            if len(parts) < 2 {
                continue
            }
            nums := strings.Split(strings.TrimSpace(parts[1]), "/")
            if len(nums) >= 2 {
                avgMs, _ = strconv.ParseFloat(nums[1], 64)
            }
        }
    }

    if avgMs == 0 && lossPct == 0 {
        return 0, 0, fmt.Errorf("could not parse ping output")
    }

    return avgMs, lossPct, nil
}

func printLatencyInfo(avgMs float64, lossPct float64) {
    fmt.Println(ui.Heading("Latency test (ping 8.8.8.8)"))

    lossColor := ui.Green
    if lossPct > 0 {
        lossColor = ui.Yellow
    }
    if lossPct >= 5 {
        lossColor = ui.Red
    }

    latencyColor := ui.Green
    switch {
    case avgMs <= 40:
    case avgMs <= 80:
        latencyColor = ui.Yellow
    default:
        latencyColor = ui.Red
    }

    fmt.Printf("  %s %s%.1f ms%s\n", ui.Key("Average:"), latencyColor, avgMs, ui.Reset)
    fmt.Printf("  %s %s%.1f%%%s\n", ui.Key("Loss   :"), lossColor, lossPct, ui.Reset)
}

func printWifiSuggestions(s wifiStatus, avgMs float64, lossPct float64) {
    fmt.Println(ui.Heading("Suggestions"))

    lines := sysinfo.WifiSuggestions(s.SignalPercent, s.Band, s.Channel, s.SecurityRaw, avgMs, lossPct)
    for _, line := range lines {
        if line == "" {
            fmt.Println()
        } else {
            fmt.Println("  " + line)
        }
    }
}

