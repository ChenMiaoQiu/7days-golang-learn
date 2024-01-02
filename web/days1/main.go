package main

import (
	"fmt"
	"net/http"

	"github.com/ChenMiaoQiu/7days-golang-learn/web/days1/gee"
)

func main() {
	c := gee.New()

	c.Get("/", indexHandler)
	c.Get("/", helloHandler)

	c.Run(":9999")
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
