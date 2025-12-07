package sysinfo

import (
   "testing"
   "strings"
)

func TestSignalQualityLabel(t *testing.T) {
    cases := []struct {
        name    string
        percent int
        want    string
    }{
        {"below zero", -5, "no signal"},
        {"zero", 0, "no signal"},
        {"very weak low", 1, "very weak"},
        {"very weak upper", 24, "very weak"},
        {"weak low", 25, "weak"},
        {"weak upper", 49, "weak"},
        {"ok low", 50, "ok"},
        {"ok upper", 74, "ok"},
        {"strong low", 75, "strong"},
        {"strong high", 100, "strong"},
        {"above hundred", 110, "strong"},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got := SignalQualityLabel(tc.percent)
            if got != tc.want {
                t.Fatalf("SignalQualityLabel(%d) = %q, want %q", tc.percent, got, tc.want)
            }
        })
    }
}

func TestBandFromFreq(t *testing.T) {
    cases := []struct {
        input string
        wantB string
        wantF int
    }{
        {"2412", "2.4 GHz band", 2412},
        {"2412 MHz", "2.4 GHz band", 2412},
        {"5180", "5 GHz band", 5180},
        {"6000", "6 GHz band", 6000},
        {"", "", 0},
        {"abc", "", 0},
    }

    for _, tc := range cases {
        gotB, gotF := BandFromFreq(tc.input)
        if gotB != tc.wantB || gotF != tc.wantF {
            t.Fatalf("BandFromFreq(%q) = (%q, %d), want (%q, %d)", tc.input, gotB, gotF, tc.wantB, tc.wantF)
        }
    }
}

func TestChannelFromFreq(t *testing.T) {
    cases := []struct {
        freq int
        want int
    }{
        {2412, 1},
        {2437, 6},
        {2462, 11},
        {2484, 14},
        {5180, (5180 - 5000) / 5},
        {0, 0},
    }

    for _, tc := range cases {
        got := ChannelFromFreq(tc.freq)
        if got != tc.want {
            t.Fatalf("ChannelFromFreq(%d) = %d, want %d", tc.freq, got, tc.want)
        }
    }
}

func TestChannelHintText(t *testing.T) {
    cases := []struct {
        band string
        ch   int
        want string
    }{
        {"2.4 GHz band", 1, "(non overlapping channel)"},
        {"2.4 GHz band", 6, "(non overlapping channel)"},
        {"2.4 GHz band", 11, "(non overlapping channel)"},
        {"2.4 GHz band", 3, "(may overlap, 1 or 6 or 11 is often better)"},
        {"5 GHz band", 36, "(usually faster, shorter range)"},
        {"6 GHz band", 5, "(usually faster, shorter range)"},
        {"", 0, ""},
    }

    for _, tc := range cases {
        got := ChannelHintText(tc.band, tc.ch)
        if got != tc.want {
            t.Fatalf("ChannelHintText(%q, %d) = %q, want %q", tc.band, tc.ch, got, tc.want)
        }
    }
}

func TestQualityToPercent(t *testing.T) {
    cases := []struct {
        in   string
        want int
    }{
        {"40/80", 50},
        {"35/70", 50},
        {"0/70", 0},
        {"70/70", 100},
        {"bad", 0},
        {"10/", 0},
        {"/10", 0},
        {"", 0},
    }

    for _, tc := range cases {
        t.Run(tc.in, func(t *testing.T) {
            got := QualityToPercent(tc.in)
            if got != tc.want {
                t.Fatalf("QualityToPercent(%q) = %d, want %d", tc.in, got, tc.want)
            }
        })
    }
}

func TestWifiSuggestionsStrongSignal(t *testing.T) {
    lines := WifiSuggestions(80, "5 GHz band", 36, "WPA2", 30, 0)

    if len(lines) == 0 {
        t.Fatalf("expected some suggestions, got none")
    }

    foundSignal := false
    for _, l := range lines {
        if strings.Contains(l, "Signal is strong") {
            foundSignal = true
            break
        }
    }
    if !foundSignal {
        t.Fatalf("expected strong signal suggestion, got: %#v", lines)
    }
}

func TestWifiSuggestionsWeakSignalAndLoss(t *testing.T) {
    lines := WifiSuggestions(10, "2.4 GHz band", 3, "WEP", 120, 10)

    hasWeak := false
    hasWep := false
    hasLoss := false

    for _, l := range lines {
        if strings.Contains(l, "Signal is weak") {
            hasWeak = true
        }
        if strings.Contains(l, "WEP") {
            hasWep = true
        }
        if strings.Contains(l, "Packet loss is noticeable") {
            hasLoss = true
        }
    }

    if !hasWeak || !hasWep || !hasLoss {
        t.Fatalf("missing expected suggestions, lines: %#v", lines)
    }
}

