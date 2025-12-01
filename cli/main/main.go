package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/tomiok/dishwasher/cli"
	"github.com/tomiok/dishwasher/internal/server"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const (
	profileServer = "server"
	profileClient = "client"
)

func main() {
	profile := os.Args[1]
	if profile == "" {
		profile = profileServer
	}

	addr := flag.String("addr", ":7000", "")
	join := flag.String("join", "", "seed node to join")
	seed := flag.Bool("seed", false, "run as seed node")
	flag.Parse()

	switch profile {
	case profileServer:
	case profileClient:
		err := runREPL(*addr)
		if err != nil {
			fmt.Printf("error exiting the client: %v\n", err)
		}
	}

	serv := server.New(*addr, *seed)

	if *seed {
		log.Fatal(serv.Start())
	}

	log.Fatalf("do not have join feature yet, %s", *join)
}

func runREPL(addr string) error {
	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("DHWS>")
LOOP:
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			fmt.Print("DHWS>")
			continue
		}

		cmd := fields[0]

		switch cmd {
		case cli.CMDPing:
			_, err = conn.Write([]byte{server.MsgPing, 0, 0, 0, 0})
			if err != nil {
				break
			}

			res := make([]byte, 5)
			_, err = io.ReadFull(conn, res)
			if err != nil {
				break LOOP
			}

			if res[0] == server.MsgPong {
				fmt.Print("DSWH>PONG\n")
			}

			fmt.Print("DSWH>")
		}
	}

	return conn.Close()
}
