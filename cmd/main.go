package main

import (
	"fmt"
	"thermalkv/internal/store"
	"time"
)

func main() {
	db := store.NewStore()

	db.Set("token", "abc123")

	db.SetTTL("token", 5)

	value, exists := db.Get("token")
	fmt.Println(value, exists)

	time.Sleep(6 * time.Second)

	value, exists = db.Get("token")
	fmt.Println(value, exists)
}
