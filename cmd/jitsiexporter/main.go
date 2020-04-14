package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/xsteadfastx/jitsiexporter"
)

func main() {
	version := "development"
	ver := flag.Bool("version", false, "Prints version.")
	url := flag.String("url", "", "URL of Jitsi Videobridge Colibri Stats.")
	debug := flag.Bool("debug", false, "Enable debug.")
	interval := flag.Duration("interval", 30*time.Second, "Seconds to wait before scraping.") // nolint: gomnd
	port := flag.Int("port", 9700, "Port to listen on.")
	host := flag.String("host", "localhost", "Host to listen on.")
	servername := flag.String("servername", "", "Jitsi server name. Used as prometheus label.")
	flag.Parse()

	if *ver {
		fmt.Print(version)
		os.Exit(0)
	}

	if *url == "" {
		log.Fatal("needs a url!")
	}

	jitsiexporter.Serve(*url, *debug, *interval, *port, *host, *servername)
}
