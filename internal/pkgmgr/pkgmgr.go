package pkgmgr

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
    "syscall"

    "penguinguide/internal/distro"
    "penguinguide/internal/ui"
)

type Options struct {
    DryRun    bool
    AssumeYes bool
    Explain   bool
}

type Manager interface {
    UpdateAll(opts Options) error
    Install(pkgs []string, opts Options) error
    Remove(pkgs []string, opts Options) error
    Search(query string, opts Options) error
    Info(name string, opts Options) error
}

func New(d *distro.Distro) Manager {
    switch d.Family {
    case distro.FamilyDebian:
        return &aptManager{}
    case distro.FamilyRHEL:
        return &dnfManager{}
    case distro.FamilyArch:
        return &pacmanManager{}
    case distro.FamilyAlpine:
        return &apkManager{}
    default:
        return &noopManager{distroID: d.ID}
    }
}

/********** APT **********/

type aptManager struct{}

func (m *aptManager) UpdateAll(opts Options) error {
    cmd := "sudo apt update && sudo apt upgrade"
    if opts.AssumeYes {
        cmd = "sudo apt update && sudo apt upgrade -y"
    }
    return runOrPrint(cmd, opts, "Update all packages with apt")
}

func (m *aptManager) Install(pkgs []string, opts Options) error {
    args := []string{"sudo", "apt", "install"}
    if opts.AssumeYes {
        args = append(args, "-y")
    }
    args = append(args, pkgs...)
    cmd := joinCommand(args)
    return runOrPrint(cmd, opts, "Install packages with apt")
}

func (m *aptManager) Remove(pkgs []string, opts Options) error {
    args := []string{"sudo", "apt", "remove"}
    if opts.AssumeYes {
        args = append(args, "-y")
    }
    args = append(args, pkgs...)
    cmd := joinCommand(args)
    return runOrPrint(cmd, opts, "Remove packages with apt")
}

func (m *aptManager) Search(query string, opts Options) error {
    cmd := "apt search " + query
    return runOrPrint(cmd, opts, "Search for packages with apt")
}

func (m *aptManager) Info(name string, opts Options) error {
    cmd := "apt show " + name
    return runOrPrint(cmd, opts, "Show package details with apt")
}

/********** DNF **********/

type dnfManager struct{}

func (m *dnfManager) UpdateAll(opts Options) error {
    cmd := "sudo dnf upgrade"
    if opts.AssumeYes {
        cmd = "sudo dnf upgrade -y"
    }
    return runOrPrint(cmd, opts, "Update all packages with dnf")
}

func (m *dnfManager) Install(pkgs []string, opts Options) error {
    args := []string{"sudo", "dnf", "install"}
    if opts.AssumeYes {
        args = append(args, "-y")
    }
    args = append(args, pkgs...)
    cmd := joinCommand(args)
    return runOrPrint(cmd, opts, "Install packages with dnf")
}

func (m *dnfManager) Remove(pkgs []string, opts Options) error {
    args := []string{"sudo", "dnf", "remove"}
    if opts.AssumeYes {
        args = append(args, "-y")
    }
    args = append(args, pkgs...)
    cmd := joinCommand(args)
    return runOrPrint(cmd, opts, "Remove packages with dnf")
}

func (m *dnfManager) Search(query string, opts Options) error {
    cmd := "dnf search " + query
    return runOrPrint(cmd, opts, "Search for packages with dnf")
}

func (m *dnfManager) Info(name string, opts Options) error {
    cmd := "dnf info " + name
    return runOrPrint(cmd, opts, "Show package details with dnf")
}

/********** Pacman **********/

type pacmanManager struct{}

func (m *pacmanManager) UpdateAll(opts Options) error {
    args := []string{"sudo", "pacman", "-Syu"}
    if opts.AssumeYes {
        args = append(args, "--noconfirm")
    }
    cmd := joinCommand(args)
    return runOrPrint(cmd, opts, "Update all packages with pacman")
}

func (m *pacmanManager) Install(pkgs []string, opts Options) error {
    args := []string{"sudo", "pacman", "-S"}
    if opts.AssumeYes {
        args = append(args, "--noconfirm")
    }
    args = append(args, pkgs...)
    cmd := joinCommand(args)
    return runOrPrint(cmd, opts, "Install packages with pacman")
}

func (m *pacmanManager) Remove(pkgs []string, opts Options) error {
    args := []string{"sudo", "pacman", "-R"}
    if opts.AssumeYes {
        args = append(args, "--noconfirm")
    }
    args = append(args, pkgs...)
    cmd := joinCommand(args)
    return runOrPrint(cmd, opts, "Remove packages with pacman")
}

func (m *pacmanManager) Search(query string, opts Options) error {
    cmd := "pacman -Ss " + query
    return runOrPrint(cmd, opts, "Search for packages with pacman")
}

func (m *pacmanManager) Info(name string, opts Options) error {
    cmd := "pacman -Si " + name
    return runOrPrint(cmd, opts, "Show package details with pacman")
}

/********** APK **********/

type apkManager struct{}

func (m *apkManager) UpdateAll(opts Options) error {
    cmd := "sudo apk update && sudo apk upgrade"
    return runOrPrint(cmd, opts, "Update all packages with apk")
}

func (m *apkManager) Install(pkgs []string, opts Options) error {
    args := []string{"sudo", "apk", "add"}
    args = append(args, pkgs...)
    cmd := joinCommand(args)
    return runOrPrint(cmd, opts, "Install packages with apk")
}

func (m *apkManager) Remove(pkgs []string, opts Options) error {
    args := []string{"sudo", "apk", "del"}
    args = append(args, pkgs...)
    cmd := joinCommand(args)
    return runOrPrint(cmd, opts, "Remove packages with apk")
}

func (m *apkManager) Search(query string, opts Options) error {
    cmd := "apk search " + query
    return runOrPrint(cmd, opts, "Search for packages with apk")
}

func (m *apkManager) Info(name string, opts Options) error {
    cmd := "apk info -a " + name
    return runOrPrint(cmd, opts, "Show package details with apk")
}

/********** Fallback **********/

type noopManager struct {
    distroID string
}

func (m *noopManager) UpdateAll(opts Options) error {
    return fmt.Errorf("package manager not implemented for distro %q", m.distroID)
}

func (m *noopManager) Install(pkgs []string, opts Options) error {
    return fmt.Errorf("install not implemented for distro %q", m.distroID)
}

func (m *noopManager) Remove(pkgs []string, opts Options) error {
    return fmt.Errorf("remove not implemented for distro %q", m.distroID)
}

func (m *noopManager) Search(query string, opts Options) error {
    return fmt.Errorf("search not implemented for distro %q", m.distroID)
}

func (m *noopManager) Info(name string, opts Options) error {
    return fmt.Errorf("info not implemented for distro %q", m.distroID)
}

/********** Helpers **********/

func runOrPrint(command string, opts Options, explanation string) error {
    if opts.Explain {
        fmt.Println(ui.Heading("Explanation"))
        if explanation != "" {
            fmt.Println("  " + explanation)
        }
        fmt.Println()
        fmt.Println(ui.Key("Native command:"))
        fmt.Println("  " + ui.Value(command))
        fmt.Println()
    }

    if opts.DryRun {
        fmt.Println(ui.Heading("Dry run"))
        fmt.Println("  Would run:")
        fmt.Println("    " + ui.Value(command))
        fmt.Println()
        fmt.Print(ui.Info("Do you want to run this command now") + " [y/N]: ")

        var response string
        _, err := fmt.Fscan(os.Stdin, &response)
        if err != nil {
            fmt.Println()
            fmt.Println(ui.Muted("Skipped running command."))
            return nil
        }

        response = strings.TrimSpace(strings.ToLower(response))
        if response != "y" && response != "yes" {
            fmt.Println(ui.Muted("Skipped running command."))
            return nil
        }

        opts.DryRun = false
        fmt.Println()
    }

    if !opts.Explain && explanation != "" {
        fmt.Println(ui.Info(explanation))
    }

    fmt.Println(ui.Key("Running:"))
    fmt.Println("  " + ui.Value(command))

    cmd := exec.Command("sh", "-c", command)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    err := cmd.Run()
    if err == nil {
        return nil
    }

    if exitErr, ok := err.(*exec.ExitError); ok {
        if status, ok2 := exitErr.Sys().(syscall.WaitStatus); ok2 {
            if status.Signaled() && status.Signal() == syscall.SIGINT {
                fmt.Println()
                fmt.Println(ui.Muted("Command canceled by user."))
                return nil
            }
        }
    }

    return err
}

func joinCommand(parts []string) string {
    return strings.Join(parts, " ")
}

