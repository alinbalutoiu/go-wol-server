package main

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	port     int
	logLevel string
	wait     time.Duration
	macs     string
)

func initLog() {
	logrusLogLevel, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("Invalid log level: %v (see --help for more info)", logLevel)
	}

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	log.SetLevel(logrusLogLevel)

	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
}

func main() {
	app := cli.NewApp()
	app.Name = "Golang Wake-on-LAN Server"
	app.Usage = "Starts the WOL server"
	app.Authors = []*cli.Author{
		&cli.Author{
			Name: "Alin Balutoiu",
		},
	}
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "loglevel",
			Usage:       "sets the log level (error, warn, info, debug)",
			Value:       "info",
			EnvVars:     []string{"LOG_LEVEL"},
			Destination: &logLevel,
		},
		&cli.IntFlag{
			Name:        "port",
			Usage:       "port to listen to",
			EnvVars:     []string{"PORT"},
			Value:       8080,
			Destination: &port,
		},
		&cli.DurationFlag{
			Name: "graceful-timeout",
			Usage: ("the duration for which the server gracefully wait for " +
				"existing connections to finish - e.g. 15s or 1m"),
			EnvVars:     []string{"GRACEFUL_TIMEOUT"},
			Value:       15 * time.Second,
			Destination: &wait,
		},
	}

	app.Action = func(c *cli.Context) error {
		return runApp(c)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runApp(c *cli.Context) error {
	initLog()

	a := NewApp(port, wait)

	a.Initialize()
	log.Infof("App initialized")

	// Initialize app router
	a.InitializeRouter()
	log.Infof("Gorilla router initialized")

	// Start serving requests to signal that the app is healthy (but not
	// necessarily ready to serve requests yet)
	log.Infof("Listening on port: %v...", port)
	return a.Run()
}
