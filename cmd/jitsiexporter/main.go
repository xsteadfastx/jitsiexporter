// nolint: gochecknoglobals
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"git.xsfx.dev/prometheus/jitsiexporter"
	log "github.com/sirupsen/logrus"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	ver := flag.Bool("version", false, "Prints version.")
	url := flag.String("url", "", "URL of Jitsi Videobridge Colibri Stats.")
	debug := flag.Bool("debug", false, "Enable debug.")
	interval := flag.Duration("interval", 30*time.Second, "Seconds to wait before scraping.") // nolint: gomnd
	port := flag.Int("port", 9700, "Port to listen on.")
	host := flag.String("host", "localhost", "Host to listen on.")
	flag.Parse()

	if *ver {
		fmt.Printf("jitsiexporter %s, commit %s, built at %s", version, commit, date)
		os.Exit(0)
	}

	if *url == "" {
		log.Fatal("needs a url!")
	}

	jitsiexporter.Serve(*url, *debug, *interval, *port, *host)
}
