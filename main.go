package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

const usage = `puppet-agent-exporter - prometheus exporter for puppet agent

Usage:

  puppet-agent-exporter [commands|flags]

The commands & flags are:
`

func main() {
	var (
		puppetYamlSummaryReportFile = flag.String("puppet.last-run-summary", "/var/lib/puppet/state/last_run_summary.yaml", "Path to puppet last_run_summary.yaml file")
		puppetYamlFullReportFile    = flag.String("puppet.last-run-report", "", "Path to puppet last_run_report.yaml file")
		puppetDisabledLockFile      = flag.String("puppet.disabled-lock", "", "Path to puppet agent_disabled.lock file")
		namespace                   = flag.String("namespace", "puppet_agent", "Namespace for metrics")
		listenAddress               = flag.String("web.listen-address", ":9001", "Address to listen on for web interface and telemetry.")
		metricsPath                 = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	)

	flag.Usage = printUsage
	flag.Parse()

	handleFlags(flag.Args())

	var reportFileScraper PuppetYamlReportScraper
	if *puppetYamlFullReportFile != "" {
		reportFileScraper = NewFullReportScraper(*namespace, *puppetYamlFullReportFile, *puppetDisabledLockFile)
	} else {
		reportFileScraper = NewSummaryReportScraper(*namespace, *puppetYamlSummaryReportFile, *puppetDisabledLockFile)
	}

	prometheus.MustRegister(
		NewPuppetExporter(
			*namespace,
			reportFileScraper,
		),
	)

	serveHTTP(*listenAddress, *metricsPath)
}

func serveHTTP(listenAddress string, metricsEndpoint string) {
	http.Handle(metricsEndpoint, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Puppet Agent Exporter</title></head>
			<body>
			<h1>Puppet Agent Exporter</h1>
			<p><a href="` + metricsEndpoint + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Info("Listening on: ", formatListenAddr(listenAddress))

	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

// formatListenAddr returns formatted UNIX addr
func formatListenAddr(addr string) string {
	parts := strings.Split(addr, ":")
	if len(parts) == 2 && parts[0] == "" {
		addr = fmt.Sprintf("localhost:%s", parts[1])
	}
	return "http://" + addr
}

// handleFlags handles cli flags
func handleFlags(flags []string) {
	if len(flags) == 0 {
		return
	}

	switch flags[0] {
	case "help":
		printUsage()
	}
}

// printUsage prints exporter usage info
func printUsage() {
	fmt.Println(usage)
	flag.PrintDefaults()

	os.Exit(0)
}
