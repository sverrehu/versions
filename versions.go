package main

import (
	"os"

	"github.com/sverrehu/gotest/versions/internal/webserver"
	"github.com/sverrehu/goutils/getopt"
)

func main() {
	help := false
	port := 8086
	opts := []getopt.Option{
		{ShortName: 'h', LongName: "help", Type: getopt.Flag, Target: &help},
		{ShortName: 'p', LongName: "port", Type: getopt.Integer, Target: &port},
	}
	getopt.Parse(&os.Args, opts, false)
	err := webserver.Run(port)
	if err != nil {
		panic(err)
	}
}
