package main

import "github.com/nemethk/claudeenv/cmd"

var version = "dev"

func main() {
	cmd.Execute(version)
}
