package cache

import "fmt"

type Cache struct {
	Id string
}

func NewCache() *Cache {
	fmt.Println("This is working...")
	return &Cache{Id: "Some ID"}
}
