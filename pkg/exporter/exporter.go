package exporter

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/44smkn/zenhub_exporter/pkg/model"
	"github.com/44smkn/zenhub_exporter/pkg/zenhub"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "zenhub"
)

var (
	BoardIssueInfo                 = prometheus.NewDesc(prometheus.BuildFQName(namespace, "board", "issue_info"), "Information about issue managed by board", []string{"workspace", "issue_number", "repository_id", "pipeline", "is_epic"}, nil)
	BoardIssueEstimatedStoryPoints = prometheus.NewDesc(prometheus.BuildFQName(namespace, "board", "issue_estimated_story_points"), "Estimated story point of each issue", []string{"workspace", "issue_number", "repository_id"}, nil)
)

// Exporter collects ZenHub stats from ZenHub API Response and exports them using
// the prometheus metrics package.
type Exporter struct {
	mutex  sync.RWMutex
	Logger log.Logger

	zenhub zenhub.Client
}

func NewExporter(zenhub zenhub.Client, logger log.Logger) *Exporter {
	return &Exporter{
		zenhub: zenhub,
		Logger: logger,
	}
}

// Describe describes all the metrics ever exported by the ZenHub exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- BoardIssueInfo
	ch <- BoardIssueEstimatedStoryPoints
}

// Collect fetches the stats from ZenHub API Responce and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	ctx := context.Background()
	issues, err := e.GetWorkspaceIssues(ctx)
	if err != nil {
		level.Error(e.Logger).Log("msg", err.Error())
		return
	}

	for _, is := range issues {
		ch <- prometheus.MustNewConstMetric(BoardIssueInfo, prometheus.GaugeValue, 1, is.WorkspaceName, is.IssueNumber, is.RepoID, is.PipelineName, strconv.FormatBool(is.IsEpic))
		if is.Estimate != nil {
			ch <- prometheus.MustNewConstMetric(BoardIssueEstimatedStoryPoints, prometheus.GaugeValue, float64(*is.Estimate), is.WorkspaceName, is.IssueNumber, is.RepoID)
		}
	}
}

func (e *Exporter) GetWorkspaceIssues(ctx context.Context) ([]model.Issue, error) {
	issues, err := e.zenhub.FetchWorkspaceIssues(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workspace issues: %w", err)
	}
	return issues, nil
}
