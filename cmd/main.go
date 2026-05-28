package main

import (
	// "bufio"
	// "fmt"
	// "os"
	// "strconv"
	// "strings"
	// "thermalkv/internal/persistence"
	"thermalkv/internal/server"
	"thermalkv/internal/store"
	//"text/scanner"
	// "sync"
)

func main() {

	db := store.NewStore()
	db.StartCleaner()

	server.Start(db)

	// snapshot := persistence.LoadSnapshot()
	// db.ImportData(snapshot)
	// logs := persistence.LoadLogs()
	// db.Recover(logs)

	// scanner := bufio.NewScanner(os.Stdin)

	// fmt.Println("KV Started...")

	// for {
	// 	fmt.Print("> ")
	// 	scanner.Scan()

	// 	input := scanner.Text()
	// 	parts := strings.Split(input, " ")
	// 	command := strings.ToUpper(parts[0])

	// 	switch command {

	// 	case "SET":
	// 		if len(parts) < 3 {
	// 			fmt.Println("Usage: SET key value")
	// 			continue
	// 		}
	// 		key := parts[1]
	// 		value := parts[2]
	// 		db.Set(key, value)
	// 		fmt.Println("Done :)")

	// 	case "GET":
	// 		if len(parts) < 2 {
	// 			fmt.Println("Usage: GET key")
	// 			continue
	// 		}
	// 		key := parts[1]
	// 		value, exists := db.Get(key)

	// 		if exists {
	// 			fmt.Println(value)
	// 		} else {
	// 			fmt.Println("Key not found... :(")
	// 		}

	// 	case "DEL":
	// 		if len(parts) < 2 {
	// 			fmt.Println("Usage: DEL key")
	// 			continue
	// 		}
	// 		key := parts[1]
	// 		db.Delete(key)
	// 		fmt.Println("Deleted")

	// 	case "TTL":
	// 		if len(parts) < 3 {
	// 			fmt.Println("Usage: TTL key seconds")
	// 			continue
	// 		}

	// 		key := parts[1]
	// 		seconds, err := strconv.Atoi(parts[2])

	// 		if err != nil {
	// 			fmt.Println("Invalid seconds")
	// 			continue
	// 		}
	// 		_, exists := db.Get(key)

	// 		if exists {
	// 			db.SetTTL(key, seconds)
	// 			fmt.Println("TTL set")
	// 		} else {
	// 			fmt.Println("Key not found... :(")
	// 		}

	// 	case "EXIT":
	// 		fmt.Println("bye... ")
	// 		return

	// 	default:
	// 		fmt.Println("Unknown command")
	// 	}

	// }

}
