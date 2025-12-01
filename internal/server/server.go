package server

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"strings"
)

const (
	MsgPing byte = iota + 1
	MsgPong
	MsgJoin
	MsgMembers
)

type Server struct {
	port string
	conn net.Conn
}

func New(port string) Server {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return Server{
		port: port,
	}
}
func (s Server) Start() error {
	l, err := net.Listen("tcp4", s.port)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("cannot accept conn %v", err)
			continue
		}
		s.conn = conn
		go s.HandleConn(conn)
	}
}

// HandleConn
// [1 byte: tipo de mensaje][4 bytes: largo del payload][N bytes: payload]
func (s Server) HandleConn(conn net.Conn) {
LOOP:
	for {
		format := make([]byte, 1)
		_, err := io.ReadFull(conn, format)
		if err != nil {
			_ = s.Close()
			break
		}
		msgFormat := format[0]
		msgLenBuf := make([]byte, 4)
		_, err = io.ReadFull(conn, msgLenBuf)
		msgLen := binary.LittleEndian.Uint32(msgLenBuf)

		messageBuf := make([]byte, msgLen)
		if _, err = io.ReadFull(conn, messageBuf); err != nil {
			_ = s.Close()
			break LOOP
		}

		if err = s.handleMessage(conn, msgFormat); err != nil {
			_ = s.Close()
			break LOOP
		}
	}
}

func (s Server) handleMessage(conn net.Conn, msgFormat byte) error {
	switch msgFormat {
	case MsgPing:
		if _, err := conn.Write([]byte{MsgPong, 0, 0, 0, 0}); err != nil {
			return err
		}
	case MsgJoin:
		log.Printf("JOIN \n")
	case MsgMembers:
		log.Printf("members")
	default:
		log.Print("unknown header\n")
	}
	return nil
}

func handleJoin() {

}

func (s Server) Close() error {
	return s.conn.Close()
}
