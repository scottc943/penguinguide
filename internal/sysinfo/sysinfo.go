package sysinfo

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "time"

    "penguinguide/internal/distro"
)

type SystemSummary struct {
    Hostname     string
    DistroName   string
    Kernel       string
    Uptime       string
    LoadAverage  string
    MemoryPretty string
}

// GetSystemSummary gathers basic system information for display.
func GetSystemSummary() (*SystemSummary, error) {
    // Hostname
    hostname, err := os.Hostname()
    if err != nil {
        hostname = "unknown"
    }

    // Distro
    d, err := distro.Detect()
    distroName := "unknown"
    if err == nil {
        if d.PrettyName != "" {
            distroName = d.PrettyName
        } else if d.Name != "" {
            distroName = d.Name
        } else if d.ID != "" {
            distroName = d.ID
        }
    }

    // Kernel
    kernel := "unknown"
    if out, err := exec.Command("uname", "-r").Output(); err == nil {
        kernel = strings.TrimSpace(string(out))
    }

    uptimeStr := readUptime()
    loadStr := readLoadAvg()
    memPretty := readMemInfoPretty()

    return &SystemSummary{
        Hostname:     hostname,
        DistroName:   distroName,
        Kernel:       kernel,
        Uptime:       uptimeStr,
        LoadAverage:  loadStr,
        MemoryPretty: memPretty,
    }, nil
}

func readUptime() string {
    f, err := os.Open("/proc/uptime")
    if err != nil {
        return "unknown"
    }
    defer f.Close()

    var secondsFloat float64
    _, err = fmt.Fscan(f, &secondsFloat)
    if err != nil {
        return "unknown"
    }

    seconds := int64(secondsFloat)
    d := time.Duration(seconds) * time.Second

    days := d / (24 * time.Hour)
    d -= days * 24 * time.Hour
    hours := d / time.Hour
    d -= hours * time.Hour
    minutes := d / time.Minute

    parts := []string{}
    if days > 0 {
        parts = append(parts, fmt.Sprintf("%d days", days))
    }
    if hours > 0 {
        parts = append(parts, fmt.Sprintf("%d hours", hours))
    }
    parts = append(parts, fmt.Sprintf("%d minutes", minutes))

    return strings.Join(parts, " ")
}

func readLoadAvg() string {
    f, err := os.Open("/proc/loadavg")
    if err != nil {
        return "unknown"
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    if !scanner.Scan() {
        return "unknown"
    }

    fields := strings.Fields(scanner.Text())
    if len(fields) < 3 {
        return "unknown"
    }

    return fmt.Sprintf("%s %s %s", fields[0], fields[1], fields[2])
}

func readMemInfoPretty() string {
    f, err := os.Open("/proc/meminfo")
    if err != nil {
        return "unknown"
    }
    defer f.Close()

    var memTotalKB int64
    var memAvailableKB int64

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "MemTotal:") {
            memTotalKB = parseMemKB(line)
        } else if strings.HasPrefix(line, "MemAvailable:") {
            memAvailableKB = parseMemKB(line)
        }
    }

    if memTotalKB <= 0 {
        return "unknown"
    }

    usedKB := memTotalKB - memAvailableKB
    usedGiB := float64(usedKB) / 1024.0 / 1024.0
    totalGiB := float64(memTotalKB) / 1024.0 / 1024.0
    percent := float64(usedKB) / float64(memTotalKB) * 100

    return fmt.Sprintf("%.1f GiB / %.1f GiB (%.0f%%)", usedGiB, totalGiB, percent)
}

func parseMemKB(line string) int64 {
    fields := strings.Fields(line)
    if len(fields) < 2 {
        return 0
    }
    val, err := strconv.ParseInt(fields[1], 10, 64)
    if err != nil {
        return 0
    }
    return val
}

