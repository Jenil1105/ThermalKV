package main

import (
	"fmt"
	"thermalkv/internal/store"
)

func main() {
	db := store.NewStore()

	fmt.Println("ThermalKV starting...")
	fmt.Println(db.Data)
}
