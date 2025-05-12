package main

import (
	"github.com/yeisme/notevault/cmd"
)

const (
	version     string = "v0.0.1"
	ENVIRONMENT string = "Debug"
)

func main() {
	cmd.Execute(version, ENVIRONMENT)
}
