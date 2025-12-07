# Contributing to Penguinguide

Thank you for taking the time to look at this project. Penguinguide is meant to help new Linux users learn and feel more comfortable with their system. Contributions that improve clarity, teaching, or usability are very welcome.

This document explains how the project is laid out and how you can add new commands or features.

---

## Project layout

The main directories look like this:

    cmd/               CLI command files that wire everything together
    internal/ui/       Color and formatting helpers
    internal/sysinfo/  System and network helpers
    internal/pkgmgr/   Package manager detection and actions

Commands live under `cmd/` and call into helpers under `internal/`.

---

## Adding a new command

To add a new command, create a file in `cmd/`. For example:

    cmd/example.go

A minimal command looks like this:

```go
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"

    "penguinguide/internal/ui"
)

var exampleCmd = &cobra.Command{
    Use:   "example",
    Short: "Short summary shown in help",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println(ui.Heading("Example command"))
        fmt.Println("This is where your logic goes.")
    },
}

func init() {
    RootCmd.AddCommand(exampleCmd)
}
```

---

Guidelines:

    Use functions from internal/ui for headings, labels, and colors

    Keep the Run function focused on input and output

    Move longer logic into internal/ packages when it grows

    Think from a beginner point of view when you print messages


After adding a command, you can run:

```go
go run .
```
and check that the command appears under penguinguide help.

---

## Working with system helpers

Code that gathers system and network information lives under: internal/sysinfo/

Commands such as these use it:

    cmd/sys.go

    cmd/sys_network.go

    cmd/sys_wifi.go

    cmd/sys_speedtest.go

When you add or change system features:

    Keep explanations short and clear

    Avoid unnecessary noise in the output

    Try to explain what a value means, not just show it

---

## Package manager logic

Package manager support lives in:

internal/pkgmgr/pkgmgr.go

Each manager implements methods such as:

- UpdateAll

- Install

- Remove

- Search

- Info

The New function chooses the right manager based on the detected distro family.

If you are adding support for another family or improving commands:

- Add or update a manager struct with the correct commands for that family

- Wire it into New to match the correct distro.Family

- Make sure you route everything through runOrPrint so --dry-run, --yes, and --explain work properly

---

## Style and tone

Penguinguide is written with new users in mind. When in doubt:

- Prefer plain language over jargon

- Explain what a command does when it might not be obvious

- Remember to be gentle and encouraging to new users

If a message might make someone feel stuck or overwhelmed, there is probably a softer way to phrase it.

---

## Tests and formatting

Before opening a pull request:

Format the code:

```go
go fmt ./...
```

Run basic checks:

```go
    go build ./...
    go test ./...
```

If tests are missing for a change that deserves them, feel free to add small focused tests.

## Opening pull requests

When you open a pull request:

- Describe what you changed and why

- Include sample output if your change affects user facing text

- Mention any follow up ideas if you think of future work

## Feature ideas and feedback

You can open an issue if you:

- Have an idea for a new command or helper

- See a place where the explanation could be clearer

- Find a bug or edge case

-  Want to talk through a larger change before you write code

Thanks again for helping make Linux easier and more welcoming for someone else.
