package cmd

import (
    "fmt"
    "io"
    "math/rand"
    "net/http"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "time"

    "github.com/spf13/cobra"

    "penguinguide/internal/ui"
)

var (
    speedQuick  bool
    speedSize   string
    speedUpload bool
)

var sysSpeedTestCmd = &cobra.Command{
    Use:   "speedtest",
    Short: "Run a network speed test",
    Run: func(cmd *cobra.Command, args []string) {
        runSpeedtestWithParams(true, speedQuick, speedSize, speedUpload)
    },
}

func init() {
    sysCmd.AddCommand(sysSpeedTestCmd)

    sysSpeedTestCmd.Flags().BoolVar(&speedQuick, "quick", false, "run a smaller, faster test")
    sysSpeedTestCmd.Flags().StringVar(&speedSize, "size", "", "download size, for example 25MB (default 100MB, quick uses 20MB)")
    sysSpeedTestCmd.Flags().BoolVar(&speedUpload, "upload", false, "include a simple upload test")
}

func runSpeedtestQuickNonInteractive() {
    runSpeedtestWithParams(false, true, "", false)
}

func runSpeedtestWithParams(interactive bool, quick bool, sizeStr string, upload bool) {
    if interactive {
        fmt.Print(ui.Info("Speed tests download data") + " [y/N]: ")
        var ans string
        fmt.Fscan(os.Stdin, &ans)
        ans = strings.ToLower(strings.TrimSpace(ans))
        if ans != "y" && ans != "yes" {
            fmt.Println(ui.Muted("Canceled."))
            return
        }
    } else {
        fmt.Println(ui.Heading("Quick speed test for WiFi doctor"))
    }

    if interactive && !quick && sizeStr == "" && !upload {
        if _, err := exec.LookPath("speedtest-cli"); err == nil {
            fmt.Println(ui.Info("Running full speedtest (speedtest-cli)"))
            cmd := exec.Command("speedtest-cli")
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr
            cmd.Run()
            return
        }
    }

    sizeBytes, label := chooseDownloadSize(quick, sizeStr)
    fmt.Printf("%s %s\n", ui.Key("Download size:"), ui.Value(label))

    if err := runDownloadTest(sizeBytes); err != nil {
        fmt.Println(ui.Error("Download test failed:"), err)
        return
    }

    if upload {
        fmt.Println()
        fmt.Println(ui.Info("Running basic upload test"))
        if err := runUploadTest(sizeBytes / 4); err != nil {
            fmt.Println(ui.Error("Upload test failed:"), err)
            return
        }
    }
}

func chooseDownloadSize(quick bool, sizeStr string) (int64, string) {
    mb := 100
    if quick {
        mb = 20
    }

    if sizeStr != "" {
        parsed, err := parseSizeMB(sizeStr)
        if err == nil && parsed > 0 {
            mb = parsed
        }
    }

    bytes := int64(mb) * 1024 * 1024
    label := fmt.Sprintf("%dMB", mb)
    return bytes, label
}

func parseSizeMB(s string) (int, error) {
    v := strings.TrimSpace(strings.ToUpper(s))
    v = strings.TrimSuffix(v, "MB")
    v = strings.TrimSuffix(v, "M")
    v = strings.TrimSpace(v)
    if v == "" {
        return 0, fmt.Errorf("empty size")
    }
    n, err := strconv.Atoi(v)
    if err != nil {
        return 0, err
    }
    return n, nil
}

func runDownloadTest(sizeBytes int64) error {
    mirrors := []string{
        fmt.Sprintf("https://speed.cloudflare.com/__down?bytes=%d", sizeBytes),
        "https://proof.ovh.net/files/100Mb.dat",
        "https://speedtest.reliableservers.com/100MB.test",
    }

    var resp *http.Response
    var err error

    for _, url := range mirrors {
        fmt.Println(ui.Muted("Trying mirror:"), url)
        resp, err = http.Get(url)
        if err != nil {
            fmt.Println("  ", ui.Warning("Mirror error:"), err)
            continue
        }
        if resp.StatusCode != http.StatusOK {
            fmt.Println("  ", ui.Warning("Mirror status:"), resp.Status)
            resp.Body.Close()
            continue
        }
        fmt.Println(ui.Info("Download test started"))
        break
    }

    if resp == nil || err != nil {
        if err == nil {
            err = fmt.Errorf("no mirror succeeded")
        }
        return err
    }
    defer resp.Body.Close()

    start := time.Now()
    var bytes int64

    buf := make([]byte, 8192)
    for {
        n, er := resp.Body.Read(buf)
        if n > 0 {
            bytes += int64(n)
        }
        if er == io.EOF {
            break
        }
        if er != nil {
            return er
        }
    }

    elapsed := time.Since(start).Seconds()
    if elapsed <= 0 {
        elapsed = 0.000001
    }

    mb := float64(bytes) / 1024.0 / 1024.0
    speedMBs := mb / elapsed
    speedMbit := speedMBs * 8.0

    fmt.Printf("  %s %.1f MB\n", ui.Key("Downloaded:"), mb)
    fmt.Printf("  %s %.2f seconds\n", ui.Key("Time     :" ), elapsed)
    fmt.Printf("  %s %.2f MB/s (%.2f Mbps)\n",
        ui.Key("Bandwidth:"), speedMBs, speedMbit)

    return nil
}

func runUploadTest(sizeBytes int64) error {
    if sizeBytes <= 0 {
        sizeBytes = 5 * 1024 * 1024
    }

    url := "https://httpbin.org/post"
    fmt.Printf("%s %.1f MB to %s\n",
        ui.Key("Uploading:"), float64(sizeBytes)/1024.0/1024.0, url)

    r := newRandomReader(sizeBytes)

    req, err := http.NewRequest("POST", url, r)
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/octet-stream")

    start := time.Now()
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    io.Copy(io.Discard, resp.Body)
    resp.Body.Close()

    elapsed := time.Since(start).Seconds()
    if elapsed <= 0 {
        elapsed = 0.000001
    }

    mb := float64(sizeBytes) / 1024.0 / 1024.0
    speedMBs := mb / elapsed
    speedMbit := speedMBs * 8.0

    fmt.Printf("  %s %.1f MB\n", ui.Key("Uploaded:"), mb)
    fmt.Printf("  %s %.2f seconds\n", ui.Key("Time    :"), elapsed)
    fmt.Printf("  %s %.2f MB/s (%.2f Mbps)\n",
        ui.Key("Bandwidth:"), speedMBs, speedMbit)

    return nil
}

type randomReader struct {
    remaining int64
    rnd       *rand.Rand
}

func newRandomReader(size int64) *randomReader {
    src := rand.NewSource(time.Now().UnixNano())
    return &randomReader{
        remaining: size,
        rnd:       rand.New(src),
    }
}

func (r *randomReader) Read(p []byte) (int, error) {
    if r.remaining <= 0 {
        return 0, io.EOF
    }
    if int64(len(p)) > r.remaining {
        p = p[:r.remaining]
    }
    n := len(p)
    for i := 0; i < n; i++ {
        p[i] = byte(r.rnd.Intn(256))
    }
    r.remaining -= int64(n)
    if r.remaining <= 0 {
        return n, io.EOF
    }
    return n, nil
}

