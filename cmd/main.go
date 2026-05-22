package main

import (
	"fmt"
	"thermalkv/internal/store"
)

func main() {
	db := store.NewStore()

	db.Set("name:first", "jenil")
	db.Set("empty", "")

	fmt.Println("ThermalKV starting...")
	value, exists := db.Get("empty")
	if exists {
		fmt.Println("value:", value)
	} else {
		fmt.Println("Key not found...")
	}
}
