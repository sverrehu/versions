package main

import (
	"flag"

	"github.com/sverrehu/gotest/versions/internal/webserver"
)

func main() {
	port := flag.Int("port", 8086, "web server listening port for HTTP")
	flag.Parse()
	err := webserver.Run(*port)
	if err != nil {
		panic(err)
	}
}
