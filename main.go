package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"os/exec"
	"strconv"
	"github.com/AhmetOZER/dibow/server"
	"github.com/AhmetOZER/dibow/client"
)




func main() {
	check_root()
	flag.Bool("help", false, "")
	flag.Bool("h", false, "")
	flag.Usage = func() {}
	flag.Parse()

	args := flag.Args()

	subcmd := ""
	if len(args) > 0 {
		subcmd = args[0]
		args = args[1:]
	}

	switch subcmd {
	case "server":
		diserver.Server(args);
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

func check_root() {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}
	i, err := strconv.Atoi(string(output[:len(output)-1]))
	if err != nil {
		log.Fatal(err)
	}
	if i != 0 {
			log.Fatal("This program must be run as root! (sudo)")
		}
	}
