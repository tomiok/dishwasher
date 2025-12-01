package main

import (
	"errors"
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
	args := parseServerArgs()

	if err := checkRunType(args.runType); err != nil {
		panic(err)
	}

	serv := server.New(args.addr, args.seed)
	log.Fatal(serv.Start())
}

func checkRunType(t string) error {
	t = strings.ToLower(t)
	if t != typeSerer {
		return errors.New("wrong type")
	}

	return nil
}

func parseServerArgs() Args {
	if len(os.Args) < 2 {
		log.Printf("should put at least 1 arg for running the Dishwaser server")
		os.Exit(1)
	}

	runType := os.Args[1]
	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	addr := serverCmd.String("addr", ":7000", "listen address")
	seed := serverCmd.Bool("seed", false, "run as seed")
	join := serverCmd.String("join", "", "seed to join")

	if err := serverCmd.Parse(os.Args[2:]); err != nil {
		panic(err)
	}

	return Args{
		runType: runType,
		addr:    *addr,
		join:    *join,
		seed:    *seed,
	}
}
