package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"github.com/alecthomas/kong"
	"github.com/sirupsen/logrus"

	"github.com/esanchezm/mylar3_exporter/exporter"
)

var (
	version string
)

type GlobalFlags struct {
	URI       *url.URL `name:"mylar3.api-uri" help:"Mylar3 API URL" env:"MYLAR3_API_URI" required:"True"`
	APIKey    string   `name:"mylar3.api-key" help:"Mylar3 API key" env:"MYLAR3_API_KEY" required:"True"`
	Timeout   int      `name:"mylar3.timeout" help:"Timeout in seconds to connect to the server" env:"MYLAR3_TIMEOUT" default:"10"`
	VerifySSL bool     `name:"mylar3.verify-ssl" help:"Whether to verify the SSL certificate when connecting to the Mylar3 server." negatable:"True" env:"MYLAR3_VERIFY_SSL" default:"true"`

	WebListenAddress string `name:"web.listen-address" help:"Address where the exporter will listen for connections" env:"EXPORTER_LISTEN_ADDRESS" default:":9091"`
	WebTelemetryPath string `name:"web.telemetry-path" help:"Metrics expose path" env:"EXPORTER_LISTEN_PATH" default:"/metrics"`
	LogLevel         string `name:"log.level" help:"Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]" enum:"debug,info,warn,error,fatal" env:"EXPORTER_LOG_LEVEL" default:"error"`

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

	log := logrus.New()
	levels := map[string]logrus.Level{
		"debug": logrus.DebugLevel,
		"error": logrus.ErrorLevel,
		"fatal": logrus.FatalLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
	}
	log.SetLevel(levels[opts.LogLevel])

	mylar3_opts := &exporter.Mylar3Opts{
		URI:     opts.URI.String(),
		APIKey:  opts.APIKey,
		Timeout: opts.Timeout,
	}

	if !opts.VerifySSL {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	exporter := exporter.New(mylar3_opts, log)

	http.Handle(opts.WebTelemetryPath, exporter.Handler())
	log.Info("Starting exporter on ", opts.WebListenAddress)
	http.ListenAndServe(opts.WebListenAddress, nil)
}
