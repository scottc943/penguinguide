package ui

// Core ANSI color constants
const (
    Reset  = "\033[0m"
    Bold   = "\033[1m"
    Dim    = "\033[2m"

    // Base palette
    Cyan   = "\033[36m"
    Blue   = "\033[34m"
    Green  = "\033[32m"
    Yellow = "\033[33m"
    Red    = "\033[31m"
    Gray   = "\033[90m"
)

// Styled printer helpers

func Heading(text string) string {
    return Cyan + Bold + text + Reset
}

func Success(text string) string {
    return Green + text + Reset
}

func Warning(text string) string {
    return Yellow + text + Reset
}

func Error(text string) string {
    return Red + Bold + text + Reset
}

func Info(text string) string {
    return Blue + text + Reset
}

func Muted(text string) string {
    return Gray + text + Reset
}

func Key(text string) string { // label style
    return Bold + Cyan + text + Reset
}

func Value(text string) string { // value style
    return Bold + text + Reset
}

