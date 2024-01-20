package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ChenMiaoQiu/7days-golang-learn/Cache/day4/geecache"
)

var db = map[string]string{
	"tom":  "100",
	"jack": "200",
	"sam":  "700",
}

func main() {
	geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[slow db] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))
	addr := "localhost:9999"
	peers := geecache.NewHTTPPool(addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
