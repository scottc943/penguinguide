package distro

import (
    "bufio"
    "errors"
    "os"
    "strings"
)

type Family string

const (
    FamilyDebian Family = "debian"
    FamilyRHEL   Family = "rhel"
    FamilyArch   Family = "arch"
    FamilySUSE   Family = "suse"
    FamilyAlpine Family = "alpine"
    FamilyOther  Family = "other"
)

type Distro struct {
    ID         string
    IDLike     []string
    Name       string
    PrettyName string
    VersionID  string
    Family     Family
}

// Detect reads /etc/os-release and returns a normalized Distro description.
func Detect() (*Distro, error) {
    f, err := os.Open("/etc/os-release")
    if err != nil {
        return nil, err
    }
    defer f.Close()

    values := make(map[string]string)

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }
        key := strings.ToUpper(strings.TrimSpace(parts[0]))
        val := strings.TrimSpace(parts[1])

        // remove surrounding quotes if present
        val = strings.Trim(val, `"'`)

        values[key] = val
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }

    id := values["ID"]
    if id == "" {
        return nil, errors.New("could not find ID in /etc/os-release")
    }

    idLikeRaw := values["ID_LIKE"]
    var idLike []string
    if idLikeRaw != "" {
        fields := strings.Fields(idLikeRaw)
        for _, f := range fields {
            f = strings.TrimSpace(f)
            if f != "" {
                idLike = append(idLike, f)
            }
        }
    }

    d := &Distro{
        ID:         id,
        IDLike:     idLike,
        Name:       values["NAME"],
        PrettyName: values["PRETTY_NAME"],
        VersionID:  values["VERSION_ID"],
    }
    d.Family = classifyFamily(d)

    return d, nil
}

func classifyFamily(d *Distro) Family {
    id := strings.ToLower(d.ID)
    like := make([]string, 0, len(d.IDLike))
    for _, v := range d.IDLike {
        like = append(like, strings.ToLower(v))
    }

    hasLike := func(target string) bool {
        for _, v := range like {
            if v == target {
                return true
            }
        }
        return false
    }

    // Debian family
    if id == "debian" || id == "ubuntu" || id == "linuxmint" || id == "raspbian" ||
        hasLike("debian") {
        return FamilyDebian
    }

    // RHEL family
    if id == "rhel" || id == "centos" || id == "rocky" || id == "almalinux" || id == "fedora" ||
        hasLike("rhel") || hasLike("fedora") {
        return FamilyRHEL
    }

    // Arch family
    if id == "arch" || id == "manjaro" || hasLike("arch") {
        return FamilyArch
    }

    // SUSE family
    if strings.Contains(id, "suse") || hasLike("suse") {
        return FamilySUSE
    }

    // Alpine
    if id == "alpine" {
        return FamilyAlpine
    }

    return FamilyOther
}

