package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/tomiok/dishwasher/cli"
	"github.com/tomiok/dishwasher/internal/server"
	"io"
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
	flag.Parse()

	switch profile {
	case profileServer:
	case profileClient:
		err := runREPL(*addr)
		if err != nil {
			fmt.Printf("error exiting the client: %v\n", err)
		}
	}
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
			prompt(cli.PromptDSWH, "")
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
				prompt("", "PONG")
			}

			prompt(cli.PromptDSWH, "")
		case cli.CMDMembers:
			_, err = conn.Write([]byte{server.MsgMembers, 0, 0, 0, 0})
			if err != nil {
				break
			}

			header := make([]byte, 5)
			if _, err = io.ReadFull(conn, header); err != nil {
				return err
			}

			if header[0] != server.MsgMembers {
				return fmt.Errorf("expected MsgMembers, got %d", header[0])
			}

			length := binary.LittleEndian.Uint32(header[1:5])

			if length > 0 {
				respPayload := make([]byte, length)
				if _, err = io.ReadFull(conn, respPayload); err != nil {
					return err
				}

				prompt("", string(respPayload))
			}

			prompt(cli.PromptDSWH, "")
		}
	}

	return conn.Close()
}

func prompt(ps1, value string) {
	if ps1 == "" {
		fmt.Print(value + "\n")
		return
	}

	if value == "" {
		fmt.Print(fmt.Sprintf("%s", ps1))
		return
	}

	fmt.Print(fmt.Sprintf("%s%s\n", ps1, value))
}
