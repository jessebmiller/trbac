package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/jessebmiller/trbac/proxy"
)

func main() {
	log.Println("Trbac proxy starting up...")
	proxy := conf.ConfiguredProxy()
	http.Handle(proxy)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
