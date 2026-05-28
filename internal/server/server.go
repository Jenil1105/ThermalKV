package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"thermalkv/internal/store"
)

func Start(db *store.Store) {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Server error:", err)
		return
	}

	defer listener.Close()

	fmt.Println("ThermalKV server running on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr())
		go HandleConnection(conn, db)
	}
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
		parts := strings.SplitN(input, " ", 3)

		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])
		switch command {

		case "SET":
			if len(parts) < 3 {
				conn.Write([]byte("Usage: SET key value\n"))
				continue
			}
			key := parts[1]
			value := parts[2]
			db.Set(key, value)
			conn.Write([]byte("Done :)\n"))

		case "GET":
			if len(parts) < 2 {
				conn.Write([]byte("Usage: GET key\n"))
				continue
			}
			key := parts[1]
			value, exists := db.Get(key)

			if exists {
				conn.Write([]byte(value + "\n"))
			} else {
				conn.Write([]byte("Key not found... :(\n"))
			}

		case "DEL":
			if len(parts) < 2 {
				conn.Write([]byte("Usage: DEL key\n"))
				continue
			}
			key := parts[1]
			db.Delete(key)
			conn.Write([]byte("Deleted\n"))

		case "TTL":
			if len(parts) < 3 {
				conn.Write([]byte("Usage: TTL key seconds\n"))
				continue
			}

			key := parts[1]
			seconds, err := strconv.Atoi(parts[2])

			if err != nil {
				conn.Write([]byte("Invalid seconds\n"))
				continue
			}
			_, exists := db.Get(key)

			if exists {
				db.SetTTL(key, seconds)
				conn.Write([]byte("TTL set\n"))
			} else {
				conn.Write([]byte("Key not found... :(\n"))
			}

		case "EXIT":
			conn.Write([]byte("bye... \n"))
			return

		default:
			conn.Write([]byte("Unknown command\n"))
		}
	}
}
