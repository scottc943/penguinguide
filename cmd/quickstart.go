package cmd

import (
    "bufio"
    "fmt"
    "os"
    "strings"

    "github.com/spf13/cobra"

    "penguinguide/internal/distro"
    "penguinguide/internal/ui"
)

var quickstartScript bool

var quickstartCmd = &cobra.Command{
    Use:   "quickstart",
    Short: "Guided tour for new Linux users",
    Long: `quickstart walks through common tasks
such as viewing system details, installing packages,
and checking network and WiFi.

It uses the same commands penguinguide provides,
but in a guided menu format or as a script of native commands.`,
    Run: func(cmd *cobra.Command, args []string) {
        if quickstartScript {
            runQuickstartScript()
        } else {
            runQuickstart()
        }
    },
}

func init() {
    RootCmd.AddCommand(quickstartCmd)
    quickstartCmd.Flags().BoolVar(&quickstartScript, "script", false, "print native commands for this quickstart")
}

func runQuickstart() {
    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Println()
        fmt.Println(ui.Heading("Penguinguide quickstart"))
        fmt.Println()
        fmt.Println("Choose a topic to explore:")
        fmt.Println("  1) See a summary of this system")
        fmt.Println("  2) Learn how package install works")
        fmt.Println("  3) Try installing a package")
        fmt.Println("  4) Review network and IP information")
        fmt.Println("  5) Check WiFi and connection health")
        fmt.Println("  0) Exit quickstart")
        fmt.Println()
        fmt.Print("Enter a number and press Enter: ")

        choice, _ := reader.ReadString('\n')
        choice = strings.TrimSpace(choice)

        fmt.Println()

        switch choice {
        case "1":
            runQuickstartSystemSummary(reader)
        case "2":
            runQuickstartInstallExplain(reader)
        case "3":
            runQuickstartInstallPackage(reader)
        case "4":
            runQuickstartNetwork(reader)
        case "5":
            runQuickstartWifi(reader)
        case "0", "q", "Q", "exit":
            fmt.Println(ui.Success("Leaving quickstart"))
            return
        default:
            fmt.Println(ui.Warning("Unknown choice, please select a number from the menu"))
        }
    }
}

func runQuickstartSystemSummary(reader *bufio.Reader) {
    fmt.Println(ui.Heading("System summary"))
    fmt.Println("This shows the machine name, distribution, kernel, uptime,")
    fmt.Println("and a basic view of load and memory usage.")
    fmt.Println()

    runSysSummary()

    waitForEnter(reader)
}

func runQuickstartInstallExplain(reader *bufio.Reader) {
    fmt.Println(ui.Heading("How package install works"))
    fmt.Println("In Linux, packages are installed with the native package manager.")
    fmt.Println("Penguinguide detects your distribution family and shows the exact")
    fmt.Println("command it uses through the package manager on your system.")
    fmt.Println()
    fmt.Println("Global flags:")
    fmt.Println("  --dry-run   show and confirm commands before running them")
    fmt.Println("  --yes       assume yes for prompts from the package manager")
    fmt.Println("  --explain   show an explanation and the native command first")
    fmt.Println()
    fmt.Println("Next, you can try a real example with the next menu option.")
    fmt.Println()

    waitForEnter(reader)
}

func runQuickstartInstallPackage(reader *bufio.Reader) {
    fmt.Println(ui.Heading("Install a package"))
    fmt.Println("You can install a package by name. A common example is htop,")
    fmt.Println("which is an interactive process viewer.")
    fmt.Println()

    fmt.Print("Enter a package name to install (default: htop): ")
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)
    if name == "" {
        name = "htop"
    }

    fmt.Println()
    fmt.Printf("You chose package %s\n", ui.Value(name))
    fmt.Println("Penguinguide will now run the normal install flow.")
    fmt.Println("This respects global flags for dry run, yes, and explain.")
    fmt.Println()

    runInstall([]string{name})

    waitForEnter(reader)
}

func runQuickstartNetwork(reader *bufio.Reader) {
    fmt.Println(ui.Heading("Network and IP information"))
    fmt.Println("This shows the default route, DNS servers, and IP addresses")
    fmt.Println("for each interface on your system.")
    fmt.Println()

    runSysNetwork()

    fmt.Println()
    fmt.Println("You can also run:")
    fmt.Println("  penguinguide sys ip")
    fmt.Println("for a shorter view that focuses on interface names and addresses.")
    fmt.Println()

    waitForEnter(reader)
}

func runQuickstartWifi(reader *bufio.Reader) {
    fmt.Println(ui.Heading("WiFi and connection health"))
    fmt.Println("This checks your WiFi signal, band, and security, then can run")
    fmt.Println("a latency test to see how the connection feels in practice.")
    fmt.Println()

    wifiCheck(true)

    fmt.Println()
    fmt.Println("For a combined WiFi check and quick speed test you can run:")
    fmt.Println("  penguinguide wifi-doctor")
    fmt.Println()

    waitForEnter(reader)
}

func waitForEnter(reader *bufio.Reader) {
    fmt.Print(ui.Muted("Press Enter to return to the quickstart menu..."))
    _, _ = reader.ReadString('\n')
}

// Script mode, prints native commands only
func runQuickstartScript() {
    d, err := distro.Detect()
    if err != nil {
        fmt.Fprintln(os.Stderr, ui.Error("Could not detect distribution for script output"))
        fmt.Fprintln(os.Stderr, "  Error:", err)
        return
    }

    family := d.Family

    var updateCmd string
    var installCmd string

    switch family {
    case distro.FamilyDebian:
        updateCmd = "sudo apt update && sudo apt upgrade"
        installCmd = "sudo apt install htop"
    case distro.FamilyRHEL:
        updateCmd = "sudo dnf upgrade"
        installCmd = "sudo dnf install htop"
    case distro.FamilyArch:
        updateCmd = "sudo pacman -Syu"
        installCmd = "sudo pacman -S htop"
    case distro.FamilyAlpine:
        updateCmd = "sudo apk update && sudo apk upgrade"
        installCmd = "sudo apk add htop"
    default:
        updateCmd = "# update packages (unknown family, edit for your system)"
        installCmd = "# install htop (unknown family, edit for your system)"
    }

    fmt.Println("# Quickstart native commands generated by penguinguide")
    fmt.Printf("# Detected family: %s\n", string(family))
    if d.Name != "" {
        fmt.Printf("# Detected name  : %s\n", d.Name)
    }
    fmt.Println()

    fmt.Println("# 1. System summary")
    fmt.Println("uname -a")
    fmt.Println("cat /etc/os-release")
    fmt.Println()

    fmt.Println("# 2. Update packages for this system")
    fmt.Println(updateCmd)
    fmt.Println()

    fmt.Println("# 3. Install a helpful process viewer")
    fmt.Println(installCmd)
    fmt.Println()

    fmt.Println("# 4. Network and IP information")
    fmt.Println("# Show interfaces and addresses")
    fmt.Println("ip addr")
    fmt.Println()
    fmt.Println("# Show default route and gateways")
    fmt.Println("ip route")
    fmt.Println()
    fmt.Println("# Optional, show public IP:")
    fmt.Println("curl https://api.ipify.org")
    fmt.Println()

    fmt.Println("# 5. WiFi information and connection health")
    fmt.Println("# On systems with NetworkManager:")
    fmt.Println("nmcli dev wifi")
    fmt.Println()
    fmt.Println("# Or with wireless tools:")
    fmt.Println("iwconfig")
    fmt.Println()
    fmt.Println("# Basic latency and packet loss check:")
    fmt.Println("ping -c 4 8.8.8.8")
    fmt.Println()

    fmt.Println("# 6. Simple download speed taste using curl")
    fmt.Println("# This downloads a test file and discards it, adjust URL as needed")
    fmt.Println("time curl -o /dev/null https://speed.cloudflare.com/__down?bytes=20000000")
}

