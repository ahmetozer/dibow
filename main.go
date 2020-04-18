package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/AhmetOZER/dibow/server"
	"github.com/AhmetOZER/dibow/client"
)




func main() {
	args := flag.Args()

	subcmd := ""
	if len(args) > 0 {
		subcmd = args[0]
		args = args[1:]
	}

	switch subcmd {
	case "server":
		diserver.Server();
	case "client":
		diclient.Client(args)
	default:
		fmt.Fprintf(os.Stderr, help)
		os.Exit(1)
	}


}
var help = `
Usage: dibow [command] [--help]
Commands:
	server - runs dibow in server mode
	client - runs dibow in client mode

Read more:
	https://github.com/AhmetOZER/dibow
`
