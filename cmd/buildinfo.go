package cmd

var (
    buildVersion = "dev"
    buildCommit  = "none"
    buildDate    = "unknown"
)

func SetBuildInfo(v, c, d string) {
    buildVersion = v
    buildCommit = c
    buildDate = d
}

