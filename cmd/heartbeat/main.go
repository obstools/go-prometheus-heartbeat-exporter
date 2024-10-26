package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/obstools/go-prometheus-heartbeat-exporter/cmd/version"
	"github.com/obstools/go-prometheus-heartbeat-exporter/pkg/heartbeat"
)

var signals, logFatalf = make(chan os.Signal, 1), log.Fatalf

// Main entrypoint
func main() {
	if err := run(os.Args); err != nil {
		logFatalf("%s\n", err)
	}
}

func run(args []string, options ...flag.ErrorHandling) error {
	failureScenario := flag.ExitOnError
	if len(options) > 0 {
		failureScenario = options[0]
	}

	ver, configPath, err := attrFromCommandLine(args, failureScenario)
	if err != nil {
		return err
	}

	if ver {
		printVersionData(os.Stdout)
		return nil
	}

	server, err := heartbeat.New(configPath)
	if err != nil {
		return err
	}

	signal.Notify(signals, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	if err := server.Start(); err != nil {
		return err
	}

	<-signals

	return server.Stop()
}

func attrFromCommandLine(args []string, options ...flag.ErrorHandling) (bool, string, error) {
	failureScenario := flag.ExitOnError
	if len(options) > 0 {
		failureScenario = options[0]
	}

	flags := flag.NewFlagSet(args[0], failureScenario)
	var (
		ver    = flags.Bool("v", false, "Prints current heartbeat version")
		config = flags.String("config", "", "Configuration path")
	)
	if err := flags.Parse(args[1:]); err != nil {
		return *ver, *config, err
	}

	return *ver, *config, nil
}

// Prints to stdout current heartbeat version
func printVersionData(writer io.Writer) {
	for _, item := range [3]string{
		"heartbeat: " + version.Version,
		"commit: " + version.GitCommit,
		"built at: " + version.BuildTime,
	} {
		fmt.Fprintln(writer, item)
	}
}
