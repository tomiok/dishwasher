package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

const (
	MsgPing byte = iota + 1
	MsgPong
	MsgJoin
	MsgMembers
)

type Server struct {
	ID      string
	port    string
	conn    net.Conn
	members map[string]string // nodeID - address
	mu      sync.RWMutex
}

func New(port string, seed bool) Server {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	members := make(map[string]string)
	id := generateID(seed)
	if seed {
		members[id] = port
	}

	return Server{
		ID:      id,
		port:    port,
		members: members,
		mu:      sync.RWMutex{},
	}
}

func (s *Server) Start(seed bool, id string) error {
	l, err := net.Listen("tcp4", s.port)
	if err != nil {
		return err
	}

	logFmt := "started dishwasher as a %s with id %s"
	var typeOfServer = "node"
	if seed {
		typeOfServer = "seed"
	}

	log.Print(fmt.Sprintf(logFmt, typeOfServer, id) + "\n")

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
func (s *Server) HandleConn(conn net.Conn) {
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

		if err = s.handleMessage(conn, msgFormat, messageBuf); err != nil {
			_ = s.Close()
			break LOOP
		}
	}
}

func (s *Server) handleMessage(conn net.Conn, msgFormat byte, msg []byte) error {
	switch msgFormat {
	case MsgPing:
		if _, err := conn.Write([]byte{MsgPong, 0, 0, 0, 0}); err != nil {
			return err
		}
	case MsgJoin:
		s.handleJoin(conn, msg)
	case MsgMembers:
		if err := s.sendMembers(conn); err != nil {
			return err
		}
	default:
		log.Print("unknown header\n")
	}
	return nil
}

func (s *Server) handleJoin(conn net.Conn, msg []byte) {
	// per doc, [0]=port [1]=node id
	msgs := strings.Split(string(msg), ",")

	s.mu.Lock()
	s.members[msgs[1]] = msgs[0]
	s.mu.Unlock()
	if err := s.sendMembers(conn); err != nil {
		log.Printf("cannot send back members, %v", err)
	}
}

func (s *Server) sendMembers(conn net.Conn) error {
	s.mu.RLock()
	var members []string
	for id, port := range s.members {
		members = append(members, port+","+id)
	}
	s.mu.RUnlock()

	payload := []byte(strings.Join(members, ";"))

	buf := make([]byte, 5+len(payload))
	buf[0] = MsgMembers
	binary.LittleEndian.PutUint32(buf[1:5], uint32(len(payload)))
	copy(buf[5:], payload)

	_, err := conn.Write(buf)
	return err
}

func (s *Server) Close() error {
	return s.conn.Close()
}
