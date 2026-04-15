package main

import (
	"os"

	"github.com/sverrehu/gotest/versions/internal/config"
	"github.com/sverrehu/gotest/versions/internal/webserver"
	"github.com/sverrehu/goutils/getopt"
)

func main() {
	help := false
	port := -1
	configFile := ""
	opts := []getopt.Option{
		{ShortName: 'h', LongName: "help", Type: getopt.Flag, Target: &help},
		{ShortName: 'p', LongName: "port", Type: getopt.Integer, Target: &port},
		{ShortName: 'c', LongName: "config", Type: getopt.String, Target: &configFile},
	}
	getopt.Parse(&os.Args, opts, false)
	err := config.LoadConfig(configFile)
	if err != nil {
		panic(err)
	}
	if port > 0 {
		config.Cfg().WebServer.Port = port
	}
	err = webserver.Run()
	if err != nil {
		panic(err)
	}
}
