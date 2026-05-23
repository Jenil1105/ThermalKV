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

	for i := 0; i < 500; i++ {

		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := "key" + strconv.Itoa(i)
			value := "value" + strconv.Itoa(i)

			db.Set(key, value)
			fmt.Println("SET:", key)

			val, ok := db.Get(key)

			if ok {
				fmt.Println("GET:", key, "=", val)
			}

			if i == 3 {
				db.SetTTL(key, 2)
				fmt.Println("TTL SET:", key)
			} else {
				db.Delete(key)
				fmt.Println("DELETED:", key)
			}

		}(i)

	}
	wg.Wait()
	fmt.Println("finished")

	// db.Set("token", "abc123")
	// db.SetTTL("token", 5)

	// value, exists := db.Get("token")
	// fmt.Println(value, exists)

	// time.Sleep(6 * time.Second)

	// value, exists = db.Get("token")
	// fmt.Println(value, exists)
}
