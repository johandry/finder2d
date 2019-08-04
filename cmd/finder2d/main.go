package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/johandry/finder2d/pkg/cli"
	"github.com/johandry/finder2d/pkg/server"
)

type config struct {
	sourceFileName string
	targetFileName string
	zero           string
	one            string
	percentage     float64
	delta          int
	output         string
	port           string
}

const envPrefix = "FINDER2D"

func main() {
	// these are the default values
	opts := &config{
		zero:       " ",
		one:        "+",
		percentage: 50.0,
		delta:      1,
		output:     "json",
		port:       "8080",
	}
	opts.Init().Read()

	var err error
	if serverMode := len(opts.targetFileName) == 0; serverMode {
		err = server.Serve(opts.port, opts.sourceFileName, opts.zero, opts.one)
	} else {
		err = cli.Execute(opts.sourceFileName, opts.targetFileName, opts.zero, opts.one, opts.percentage, opts.delta, strings.ToLower(opts.output))
	}

	if err != nil {
		fmt.Printf("[ERROR] %s", err)
		os.Exit(1)
	}
}

// Init reads the configuration from environment variables and the flags
func (c *config) Init() *config {
	flag.StringVar(&c.sourceFileName, "source", getEnv("source", c.sourceFileName), "source or source matrix file (required)")
	flag.StringVar(&c.targetFileName, "target", getEnv("target", c.targetFileName), "target or target matrix file (required)")
	flag.StringVar(&c.zero, "off", getEnv("off", c.zero), "matrix character that represents a zero or off bit")
	flag.StringVar(&c.one, "on", getEnv("on", c.one), "matrix character that represents a one or on bit")
	flag.Float64Var(&c.percentage, "p", getEnvFloat("percentage", c.percentage), "matching percentage")
	flag.IntVar(&c.delta, "d", getEnvInt("delta", c.delta), "matches blurry delta, the higher it is the less blurry patterns will find")
	flag.StringVar(&c.output, "o", getEnv("output", c.output), "output format. Availabe formats are 'text' and 'json'")
	flag.StringVar(&c.port, "port", getEnv("port", c.port), "port to start the server")

	return c
}

func (c *config) Read() {
	flag.Parse()
}

func getEnv(name string, defValue string) string {
	name = strings.ToUpper(name)
	value := os.Getenv(envPrefix + "_" + name)
	if len(value) == 0 {
		value = defValue
	}
	return value
}

func getEnvFloat(name string, defVal float64) float64 {
	valStr := getEnv(name, "")
	if len(valStr) == 0 {
		return defVal
	}
	value, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return defVal
	}
	return value
}

func getEnvInt(name string, defVal int) int {
	valStr := getEnv(name, "")
	if len(valStr) == 0 {
		return defVal
	}
	value, err := strconv.Atoi(valStr)
	if err != nil {
		return defVal
	}
	return value
}
