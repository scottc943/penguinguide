package sysinfo

import (
    "net"
    "sort"
    "os"
    "os/exec"
    "strings"
)

type InterfaceInfo struct {
    Name      string
    IsUp      bool
    IsLoopback bool
    Addresses []string
}

// GetInterfaceInfo returns a summary of network interfaces and their addresses.
func GetInterfaceInfo() ([]InterfaceInfo, error) {
    ifaces, err := net.Interfaces()
    if err != nil {
        return nil, err
    }

    var results []InterfaceInfo

    for _, iface := range ifaces {
        addrs, err := iface.Addrs()
        if err != nil {
            continue
        }

        var addrStrs []string
        for _, a := range addrs {
            addrStrs = append(addrStrs, a.String())
        }

        // Skip interfaces that have no addresses at all
        if len(addrStrs) == 0 {
            continue
        }

        info := InterfaceInfo{
            Name:      iface.Name,
            IsUp:      iface.Flags&net.FlagUp != 0,
            IsLoopback: iface.Flags&net.FlagLoopback != 0,
            Addresses: addrStrs,
        }
        results = append(results, info)
    }

    // Sort by name so output is stable
    sort.Slice(results, func(i, j int) bool {
        return results[i].Name < results[j].Name
    })

    return results, nil
}

// GetDefaultRoute returns interface and gateway from `ip route` command.
func GetDefaultRoute() (iface string, gateway string) {
    out, err := exec.Command("sh", "-c", "ip route 2>/dev/null").Output()
    if err != nil {
        return "", ""
    }

    lines := strings.Split(string(out), "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "default") {
            parts := strings.Fields(line)
            for i := 0; i < len(parts)-1; i++ {
                if parts[i] == "via" {
                    gateway = parts[i+1]
                }
                if parts[i] == "dev" {
                    iface = parts[i+1]
                }
            }
            return
        }
    }
    return
}

// GetDNSServers parses /etc/resolv.conf
func GetDNSServers() []string {
    data, err := os.ReadFile("/etc/resolv.conf")
    if err != nil {
        return nil
    }

    lines := strings.Split(string(data), "\n")
    var servers []string
    for _, line := range lines {
        if strings.HasPrefix(line, "nameserver") {
            fields := strings.Fields(line)
            if len(fields) >= 2 {
                servers = append(servers, fields[1])
            }
        }
    }
    return servers
}
