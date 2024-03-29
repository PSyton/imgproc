package main // import "github.com/psyton/imgproc/cmd/imgproc"

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"

	"imgproc/internal"
)

var revision = "unknown"

func initLogger(aLevel string) {
	level, err := log.ParseLevel(aLevel)
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: false,
		TimestampFormat:  "2006-01-02T15:04:05.0000",
	})
}

func main() {
	var options internal.Options

	p := flags.NewParser(&options, flags.Default)
	if _, e := p.ParseArgs(os.Args[1:]); e != nil {
		os.Exit(1)
	}

	initLogger(options.LogLevel)

	log.Printf("imgproc %s (log level: %s)", revision, options.LogLevel)

	if err := os.MkdirAll(options.UploadLocation, os.ModePerm); err != nil {
		log.Panicf("Can't create upload folder")
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Warn("Interrupt signal")
		cancel()
	}()

	app, err := internal.NewApplication(options)
	if err != nil {
		log.Panicf("Failed to setup application, %+v", err)
	}

	app.Run(ctx)

	log.Print("imgproc stopped")
}
