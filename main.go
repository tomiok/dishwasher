package main

import (
	"flag"
	"github.com/tomiok/dishwasher/internal/server"
	"log"
	"os"
	"strings"
)

type Args struct {
	runType string
	addr    string
	join    string
	seed    bool
}

const (
	typeSerer = "server"
)

func main() {
	run()
}

func run() {
	args := checkRunType()

	serv := server.New(args.addr, args.seed)
	if args.join != "" {
		if err := serv.JoinSeedNode(args.join); err != nil {
			log.Printf("cannot join seed %v", err)
		}

	}

	log.Fatal(serv.Start(args.seed))
}

func checkRunType() Args {
	if len(os.Args) < 2 {
		log.Printf("should put at least 1 arg for running the Dishwaser server")
		os.Exit(1)
	}
	runType := os.Args[1]

	runType = strings.ToLower(runType)
	if runType == typeSerer {
		return parseServerArgs()
	}

	panic("run with a proper run type [server]")
}

func parseServerArgs() Args {
	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	addr := serverCmd.String("addr", ":7000", "listen address")
	seed := serverCmd.Bool("seed", false, "run as seed")
	join := serverCmd.String("join", "", "seed to join")

	if err := serverCmd.Parse(os.Args[2:]); err != nil {
		panic(err)
	}

	return Args{
		addr: *addr,
		join: *join,
		seed: *seed,
	}
}
