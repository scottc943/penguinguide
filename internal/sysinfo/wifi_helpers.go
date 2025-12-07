package sysinfo

import (
    "strconv"
    "strings"
)

// SignalQualityLabel returns a human friendly label
// for a wifi signal quality percentage in the range 0 to 100.
func SignalQualityLabel(percent int) string {
    switch {
    case percent <= 0:
        return "no signal"
    case percent < 25:
        return "very weak"
    case percent < 50:
        return "weak"
    case percent < 75:
        return "ok"
    case percent <= 100:
        return "strong"
    default:
        return "strong"
    }
}

// BandFromFreq parses a frequency string like "2412" or "2412 MHz"
// and returns a band label plus the numeric MHz value.
func BandFromFreq(freq string) (string, int) {
    parts := strings.Fields(freq)
    if len(parts) == 0 {
        return "", 0
    }
    fStr := parts[0]
    fVal, err := strconv.Atoi(fStr)
    if err != nil {
        return "", 0
    }

    switch {
    case fVal >= 2400 && fVal < 2500:
        return "2.4 GHz band", fVal
    case fVal >= 4900 && fVal < 5900:
        return "5 GHz band", fVal
    case fVal >= 5900 && fVal < 7125:
        return "6 GHz band", fVal
    default:
        return "", fVal
    }
}

// ChannelFromFreq maps a frequency in MHz to a channel number.
// It covers common 2.4 GHz, 5 GHz, and 6 GHz bands.
func ChannelFromFreq(freqMHz int) int {
    if freqMHz == 0 {
        return 0
    }
    if freqMHz >= 2412 && freqMHz <= 2472 {
        return (freqMHz - 2407) / 5
    }
    if freqMHz == 2484 {
        return 14
    }
    if freqMHz >= 5000 && freqMHz <= 5900 {
        return (freqMHz - 5000) / 5
    }
    if freqMHz >= 5925 && freqMHz <= 7125 {
        return (freqMHz - 5950) / 5
    }
    return 0
}

// ChannelHintText returns a plain text hint about channel choice
// without any coloring or UI decoration.
func ChannelHintText(band string, ch int) string {
    if band == "2.4 GHz band" {
        if ch == 1 || ch == 6 || ch == 11 {
            return "(non overlapping channel)"
        }
        if ch >= 2 && ch <= 13 {
            return "(may overlap, 1 or 6 or 11 is often better)"
        }
    }
    if band == "5 GHz band" || band == "6 GHz band" {
        return "(usually faster, shorter range)"
    }
    return ""
}

// QualityToPercent converts a quality string like "40/70" to an
// integer percentage 0 to 100.
func QualityToPercent(quality string) int {
    quality = strings.TrimSpace(quality)
    if quality == "" {
        return 0
    }

    if !strings.Contains(quality, "/") {
        return 0
    }

    parts := strings.SplitN(quality, "/", 2)
    if len(parts) != 2 {
        return 0
    }

    num, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
    den, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
    if err1 != nil || err2 != nil || den <= 0 {
        return 0
    }

    return int(float64(num) / float64(den) * 100.0)
}

// ExtractBetween returns the substring between start and end.
// If either marker is not found, it returns an empty string.
func ExtractBetween(s, start, end string) string {
    i := strings.Index(s, start)
    if i == -1 {
        return ""
    }
    s = s[i+len(start):]
    j := strings.Index(s, end)
    if j == -1 {
        return ""
    }
    return s[:j]
}

// ExtractAfter returns the substring after a given prefix,
// trimmed at the first space, and trimmed of outer spaces.
func ExtractAfter(s, prefix string) string {
    i := strings.Index(s, prefix)
    if i == -1 {
        return ""
    }
    s = s[i+len(prefix):]
    if j := strings.IndexAny(s, " "); j != -1 {
        s = s[:j]
    }
    return strings.TrimSpace(s)
}

// WifiSuggestions builds user facing suggestions based on
// signal strength, band, channel, security and latency.
func WifiSuggestions(signalPercent int, band string, channel int, securityRaw string, avgMs float64, lossPct float64) []string {
    var lines []string

    // Signal based suggestions
    switch {
    case signalPercent >= 70:
        lines = append(lines, "Signal is strong, no changes needed here.")
    case signalPercent >= 40:
        lines = append(lines, "Signal is fair, if you see slowdowns try moving a little closer to the router.")
    case signalPercent > 0:
        lines = append(lines, "Signal is weak, try moving closer or reducing walls and metal between you and the router.")
    default:
        lines = append(lines, "Signal value is unknown, if you have issues check your distance to the router.")
    }

    // Band and channel
    if band == "2.4 GHz band" {
        if channel == 1 || channel == 6 || channel == 11 {
            lines = append(lines,
                "You are on 2.4 GHz, which has longer range but can be crowded.",
                "Channel looks reasonable already.",
            )
        } else if channel > 0 {
            lines = append(lines,
                "You are on 2.4 GHz which is often crowded.",
                "If performance is poor, consider channels 1, 6, or 11 to reduce overlap.",
            )
        }
    }

    if band == "5 GHz band" || band == "6 GHz band" {
        lines = append(lines,
            "You are using a higher band which is usually faster but has shorter range.",
            "If you have issues far from the router, try 2.4 GHz or move closer.",
        )
    }

    // Security
    if securityRaw == "" || securityRaw == "--" {
        lines = append(lines,
            "Network appears open. Anyone nearby can join. Use WPA2 or WPA3 if possible.",
        )
    } else if strings.Contains(securityRaw, "WEP") {
        lines = append(lines,
            "Network uses WEP which is very weak. Switch to WPA2 or WPA3.",
        )
    }

    // Latency and loss
    if avgMs > 0 || lossPct > 0 {
        lines = append(lines, "")
        lines = append(lines, "Based on the latency test:")
        if lossPct >= 5 {
            lines = append(lines, "Packet loss is noticeable. This can feel like pages hanging or video stutter.")
        } else if avgMs > 80 {
            lines = append(lines, "Latency is high. Games and calls may feel laggy.")
        } else {
            lines = append(lines, "Latency and loss look reasonable for most tasks.")
        }
    }

    return lines
}

