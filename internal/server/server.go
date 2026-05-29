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

	fmt.Println("ThermalKV server running on port 8080...")

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
			info := fmt.Sprintf("Keys: %d :: Heap Entries: %d :: Writes Since Snapshot: %d", db.Count(), db.HeapSize(), db.GetWriteCount())
			WriteResponse(conn, info)

		case "EXIT":
			WriteResponse(conn, "bye... ")
			return

		default:
			WriteResponse(conn, "Unknown command")
		}
	}
}

func WriteResponse(conn net.Conn, msg string) {
	conn.Write([]byte(msg + "\n"))
}
