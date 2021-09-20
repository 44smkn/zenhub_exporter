package main

import (
	"net/http"
	"os"
	"time"

	"github.com/44smkn/zenhub_exporter/pkg/exporter"
	"github.com/44smkn/zenhub_exporter/pkg/zenhub"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// common configuration
	webConfig     = webflag.AddFlags(kingpin.CommandLine)
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9981").Envar("WEB_LISTEN_ADDRESS").String()
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").Envar("WEB_TELEMETRY_PATH").String()

	// zenhub configuraion
	zenhubAPIToken      = kingpin.Flag("zenhub.api-token", "Token needed to requests ZenHub API").Envar("ZENHUB_API_TOKEN").String()
	zenhubWorkspaceName = kingpin.Flag("zenhub.workspace-name", "your ZenHub workspace name").Envar("ZENHUB_WORKSPACE_NAME").String()
	zenhubRepoId        = kingpin.Flag("zenhub.repo-id", "repository id bound with board, See README.md").Envar("ZENHUB_REPO_ID").String()
)

func main() {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("zenhub_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promlog.New(promlogConfig)
	level.Info(logger).Log("msg", "Starting zenhub_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "context", version.BuildContext())

	level.Info(logger).Log("msg", "Create zenhub client", "workspace", *zenhubWorkspaceName, "repository", *zenhubRepoId)
	zenhubClient := zenhub.NewClient(*zenhubAPIToken, *zenhubRepoId, *zenhubWorkspaceName, time.Duration(30*time.Second))
	zenhubExporter := exporter.NewExporter(zenhubClient, logger)
	prometheus.MustRegister(zenhubExporter, version.NewCollector("zenhub_exporter"))

	level.Info(logger).Log("msg", "Listening on address", "address", listenAddress)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
<head><title>ZenHub Exporter</title></head>
<body>
<h1>ZenHub Exporter</h1>
<p><a href='` + *metricsPath + `'>Metrics</a></p>
</body>
</html>`))
	})

	srv := &http.Server{Addr: *listenAddress}
	if err := web.ListenAndServe(srv, *webConfig, logger); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
