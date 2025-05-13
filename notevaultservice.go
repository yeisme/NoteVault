package main

import (
	_ "embed"

	"github.com/yeisme/notevault/cmd"
)

const (
	ENVIRONMENT string = "Debug"
)

//go:embed version
var version string

func main() {
	cmd.Execute(version, ENVIRONMENT)
}
