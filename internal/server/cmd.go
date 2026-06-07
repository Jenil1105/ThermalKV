package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"thermalkv/internal/store"
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

func HandleConnection(conn net.Conn, db *store.Store) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected")
			return
		}
		input = strings.TrimSpace(input)
		parts := strings.Split(input, " ")

		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])
		switch command {

		case "SET":
			if len(parts) < 3 {
				WriteResponse(conn, "Usage: SET key value")
				continue
			}
			key := parts[1]
			value := strings.Join(parts[2:], " ")
			db.Set(key, value)
			WriteResponse(conn, "Done :)")

		case "GET":
			if len(parts) < 2 {
				WriteResponse(conn, "Usage: GET key")
				continue
			}
			key := parts[1]
			value, exists := db.Get(key)

			if exists {
				WriteResponse(conn, value)
			} else {
				WriteResponse(conn, "Key not found... :(")
			}

		case "DEL":
			if len(parts) < 2 {
				WriteResponse(conn, "Usage: DEL key")
				continue
			}
			key := parts[1]
			db.Delete(key)
			WriteResponse(conn, "Deleted")

		case "TTL":
			if len(parts) < 3 {
				WriteResponse(conn, "Usage: TTL key seconds")
				continue
			}

			key := parts[1]
			seconds, err := strconv.Atoi(parts[2])

			if err != nil {
				WriteResponse(conn, "Invalid seconds")
				continue
			}
			_, exists := db.Get(key)

			if exists {
				db.SetTTL(key, seconds)
				WriteResponse(conn, "TTL set")
			} else {
				WriteResponse(conn, "Key not found... :(")
			}

		case "COOL":
			if len(parts) < 2 {
				WriteResponse(conn, "Usage: COOL key")
				continue
			}

			key := parts[1]
			err := db.CoolKey(key)

			if err != nil {
				WriteResponse(conn, err.Error())
				continue
			}
			WriteResponse(conn, "Moved to COOL storage")

		case "COUNT":
			count := db.Count()

			WriteResponse(conn, fmt.Sprintf("%d", count))

		case "EXISTS":
			if len(parts) < 2 {
				WriteResponse(conn, "Usage: EXISTS key")
				continue
			}

			key := parts[1]

			if db.Exists(key) {
				WriteResponse(conn, "true")
			} else {
				WriteResponse(conn, "false")
			}

		case "KEYS":
			keys := db.Keys()

			if len(keys) == 0 {
				WriteResponse(conn, "No keys")
				continue
			}

			WriteResponse(conn, strings.Join(keys, ", "))

		case "INFO":
			WriteResponse(conn, db.GetInfo()...)

		case "EXIT":
			WriteResponse(conn, "bye... ")
			fmt.Println("Client disconnected: ", conn.RemoteAddr())
			return

		default:
			WriteResponse(conn, "Unknown command")
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
