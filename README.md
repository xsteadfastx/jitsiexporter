# jitsiexporter

[![Build Status](https://ci.xsfx.dev/api/badges/prometheus/jitsiexporter/status.svg)](https://ci.xsfx.dev/prometheus/jitsiexporter)

A Jitsi meet prometheus exporter.

        Usage of ./jitsiexporter_linux_amd64:
          -debug
                Enable debug.
          -host string
                Host to listen on. (default "localhost")
          -interval duration
                Seconds to wait before scraping. (default 30s)
          -port int
                Port to listen on. (default 9700)
          -url string
                URL of Jitsi Videobridge Colibri Stats.
          -version
                Prints version.

## Usage

For a docker based setup, you can use the docker image [quay.io/xsteadfastx/jitsiexporter](https://quay.io/repository/xsteadfastx/jitsiexporter).

1. [Enable](https://github.com/jitsi/jitsi-videobridge/blob/master/doc/statistics.md) `/colibri/stats` for the Jitsi videobridge. When you use the Jitsi docker setup use environment variable `JVB_ENABLE_APIS=rest,colibri`.
2. Be sure that the exporter and the videobridge API can communicate. In the docker Jitsi setup: Add the `jitsiexporter` to the `jitsi-meet_meet.jitsi`-network. The url would be `http://jitsi-meet_jvb_1:8080`.
3. The `-url` URL needs to be the full url: `https://videobridge/colibri/stats`.
4. A failed scrape metric is exported as `jitsi_fetch_errors`.
