package main

import (
	"flag"
	"log"
	"strings"

	"github.com/johandry/finder2d/pkg/cli"
	"github.com/johandry/finder2d/pkg/server"
)

func main() {
	var sourceFileName, targetFileName string
	var zero, one string
	var percentage float64
	var delta int
	var output string
	var serve bool
	var port string

	flag.StringVar(&sourceFileName, "source", "", "source or source matrix file (required)")
	flag.StringVar(&targetFileName, "target", "", "target or target matrix file (required)")
	flag.StringVar(&zero, "off", " ", "matrix character that represents a zero or off bit")
	flag.StringVar(&one, "on", "+", "matrix character that represents a one or on bit")
	flag.Float64Var(&percentage, "p", 50.0, "matching percentage")
	flag.IntVar(&delta, "d", 1, "matches blurry delta, the higher it is the less blurry patterns will find")
	flag.StringVar(&output, "o", "json", "output format. Availabe formats are 'text' and 'json'")
	flag.BoolVar(&serve, "server", false, "start the server")
	flag.StringVar(&port, "port", "8080", "port to start the server")
	flag.Parse()

	if len(sourceFileName) == 0 {
		log.Fatalf("source file is required. Use the flag '--source'")
	}
	if len(targetFileName) == 0 && !serve {
		log.Fatalf("target file is required. Use the flag '--target'")
	}
	switch output {
	case "", "text", "matrix", "json":
	default:
		log.Fatalf("Unknown output format %q. Available options are: 'json' and 'text' or 'matrix'", output)
	}

	if !serve {
		if err := cli.Exec(sourceFileName, targetFileName, zero, one, percentage, delta, strings.ToLower(output)); err != nil {
			log.Fatal(err)
		}
	} else {
		server.Serve(port, sourceFileName, zero, one)
	}
}
