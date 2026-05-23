package main

import (
	"fmt"
	"strconv"
	"sync"
	"thermalkv/internal/store"
)

func main() {
	db := store.NewStore()
	db.StartCleaner()

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {

		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := "key" + strconv.Itoa(i)
			value := "value" + strconv.Itoa(i)

			db.Set(key, value)
			db.Get(key)
			db.Delete(key)
		}(i)

		wg.Wait()
		fmt.Println("finished")

	}

	// db.Set("token", "abc123")
	// db.SetTTL("token", 5)

	// value, exists := db.Get("token")
	// fmt.Println(value, exists)

	// time.Sleep(6 * time.Second)

	// value, exists = db.Get("token")
	// fmt.Println(value, exists)
}
