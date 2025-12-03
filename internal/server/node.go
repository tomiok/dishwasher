package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func (s *Server) JoinSeedNode(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	nodeID := generateID("node")
	s.nodeID = nodeID
	port := s.port
	payload := []byte(fmt.Sprintf("%s,%s", port, nodeID))

	buf := make([]byte, 5+len(payload))
	buf[0] = MsgJoin
	binary.LittleEndian.PutUint32(buf[1:5], uint32(len(payload)))
	copy(buf[5:], payload)

	_, err = conn.Write(buf)
	if err != nil {
		return err
	}

	header := make([]byte, 5)
	if _, err = io.ReadFull(conn, header); err != nil {
		return err
	}

	if header[0] != MsgMembers {
		return fmt.Errorf("expected MsgMembers, got %d", header[0])
	}

	length := binary.LittleEndian.Uint32(header[1:5])
	if length > 0 {
		respPayload := make([]byte, length)
		if _, err := io.ReadFull(conn, respPayload); err != nil {
			return err
		}

		// "7001,node1;7002,node2"
		for _, m := range strings.Split(string(respPayload), ";") {
			parts := strings.Split(m, ",")
			if len(parts) == 2 {
				s.mu.Lock()
				s.members[parts[1]] = parts[0]
				s.mu.Unlock()
			}
		}
	}

	log.Printf("joined cluster with %d members", len(s.members))
	return nil
}
