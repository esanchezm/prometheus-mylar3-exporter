package main

import (
	"fmt"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/sirupsen/logrus"

	"github.com/esanchezm/mylar3_exporter/exporter"
)

var (
	version string
)

type GlobalFlags struct {
	URI     string `name:"mylar3.api-uri" help:"Mylar3 API URL" env:"MYLAR_API_URI" default:"http://localhost:8090/api"`
	APIKey  string `name:"mylar3.api-key" help:"Mylar3 API key" env:"MYLAR_API_KEY" default:""`
	Timeout int    `name:"mylar3.timeout" help:"Timeout in seconds to connect to the server" default:"10"`

	WebListenAddress string `name:"web.listen-address" help:"Address to listen on for exporter" default:":9091"`
	WebTelemetryPath string `name:"web.telemetry-path" help:"Metrics expose path" default:"/metrics"`
	LogLevel         string `name:"log.level" help:"Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]" enum:"debug,info,warn,error,fatal" default:"error"`

	Version bool `name:"version" help:"Show exporter version"`
}

func main() {
	var opts GlobalFlags

	kong.Parse(&opts,
		kong.Name("mylar3_exporter"),
		kong.Description("Mylar3 Prometheus exporter"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": version,
		})

	if opts.Version {
		fmt.Println("mylar3_exporter - Mylar3 Prometheus exporter")
		fmt.Printf("Version: %s\n", version)
		return
	}

	logger := logrus.New()
	levels := map[string]logrus.Level{
		"debug": logrus.DebugLevel,
		"error": logrus.ErrorLevel,
		"fatal": logrus.FatalLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
	}
	logger.SetLevel(levels[opts.LogLevel])

	mylar3_opts := &exporter.Mylar3Opts{
		URI:     opts.URI,
		APIKey:  opts.APIKey,
		Timeout: opts.Timeout,
	}

	exporter := exporter.New(mylar3_opts, logger)

	http.Handle(opts.WebTelemetryPath, exporter.Handler())
	http.ListenAndServe(opts.WebListenAddress, nil)
}
