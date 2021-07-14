package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = func() {
		usage :=
			`
Usage of %s:
  install
        Install crosstools into system
  update
        Update crosstools in system
	`
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "\nFlags of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	selfPtr := flag.Bool("self", false, "Take action to crosstools-installer itself instead of crosstools")

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
	}

	if *selfPtr && flag.Arg(0) == "install" {
		InstallSelf()
	}
}
