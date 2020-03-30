package main

import (
	"flag"

	"github.com/xsteadfastx/jitsiexporter"
)

func main() {
	url := flag.String("url", "", "")
	flag.Parse()
	jitsiexporter.Serve(*url)
}
