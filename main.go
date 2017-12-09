package main

import (
	"fmt"

	"github.com/VagabondDataNinjas/gizlinebot/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	fmt.Printf("\nGROOTs linebot %v, commit %v, built at %v\n", version, commit, date)
	cmd.Execute()
}
