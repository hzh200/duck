package main

import (
	// "fmt"
	"duck/kernel/server"
	"flag"
)

func main() {
	port := flag.Int("port", 9000, "backend http server port")
	flag.Parse()
	server.StartServer(*port)
}
