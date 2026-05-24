package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	//"text/scanner"

	// "strconv"
	// "sync"
	"thermalkv/internal/persistence"
	"thermalkv/internal/store"
)

func main() {

	db := store.NewStore()
	db.StartCleaner()

	logs := persistence.LoadLogs()
	db.Recover(logs)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("KV Started...")

	for {
		fmt.Print("> ")
		scanner.Scan()

		input := scanner.Text()
		parts := strings.Split(input, " ")
		command := strings.ToUpper(parts[0])

		switch command {
		case "SET":
			if len(parts) < 3 {
				fmt.Println("Usage: SET key value")
				continue
			}
			key := parts[1]
			value := parts[2]
			db.Set(key, value)
			fmt.Println("Done :)")

		case "GET":
			if len(parts) < 2 {
				fmt.Println("Usage: GET key")
				continue
			}
			key := parts[1]
			value, exists := db.Get(key)

			if exists {
				fmt.Println(value)
			} else {
				fmt.Println("Key not found... :(")
			}

		case "DEL":
			if len(parts) < 2 {
				fmt.Println("Usage: DEL key")
				continue
			}
			key := parts[1]
			db.Delete(key)
			fmt.Println("Deleted")

		case "EXIT":
			fmt.Println("bye... ")
			return

		default:
			fmt.Println("Unknown command")
		}

	}

}
