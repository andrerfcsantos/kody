package main

import (
	"kody/cmd"
	"kody/lib/config"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Execute(config.BuildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	})
}
