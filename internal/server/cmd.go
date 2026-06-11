package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const EndMarker = "__END_RESPONSE__"

func (s *Server) Start() {

	fmt.Println("ThermalKV server running on port 8080...")

	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("Accept stopped:", err)
			return
		}
		fmt.Println("Client connected:", conn.RemoteAddr())
		go HandleConnection(conn, s.DB)
	}
}

func (s *Server) Shutdown() error {

	fmt.Println("Stopping TCP listener...")
	return s.Listener.Close()

}

func HandleConnection(conn net.Conn, db KVService) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected")
			return
		}
		lines, shouldClose := ExecuteCommand(db, strings.TrimSpace(input))
		WriteResponse(conn, lines...)

		if shouldClose {
			fmt.Println("Client disconnected: ", conn.RemoteAddr())
			return
		}

	}
}

func WriteResponse(
	conn net.Conn,
	lines ...string,
) {
	for _, line := range lines {
		conn.Write([]byte(line + "\n"))
	}

	conn.Write([]byte(EndMarker + "\n"))

}
