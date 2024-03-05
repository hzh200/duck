package main

import (
	// "fmt"
	"duck/kernel/server"
	"flag"
	"fmt"
)

func main() {
	port := flag.Int("port", 9000, "backend http server port")
	mode := flag.String("mode", "proc", "kernel running mode: dev/proc")
	flag.Parse()
	devMode := *mode == "dev"

	if devMode {
		fmt.Println("Enter development mode.")
	}

	server.StartServer(*port)
}
