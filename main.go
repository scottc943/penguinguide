package main

import "penguinguide/cmd"

var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)

func main() {
    // Pass build info into the CLI layer.
    cmd.SetBuildInfo(version, commit, date)
    cmd.Execute()
}

