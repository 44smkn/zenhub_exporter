package exporter_test

import (
	"testing"

	mocks_zenhub "github.com/44smkn/zenhub_exporter/mocks/zenhub"
	"github.com/44smkn/zenhub_exporter/pkg/exporter"
	"github.com/44smkn/zenhub_exporter/pkg/model"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/promlog"
)

func Test_defaultExporter_Collect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks_zenhub.NewMockClient(ctrl)
	mockIssues := []model.Issue{
		{
			RepoID:        "358260205",
			WorkspaceName: "Test",
			IssueNumber:   "22",
			Estimate:      nil,
			IsEpic:        false,
			PipelineName:  "Backlog",
		},
	}
	m.EXPECT().FetchWorkspaceIssues(gomock.Any()).Return(mockIssues, nil)

	promlogConfig := &promlog.Config{}
	logger := promlog.New(promlogConfig)
	e := exporter.NewExporter(m, logger)
	ch := make(chan prometheus.Metric, 10)

	go e.Collect(ch)
	for m := range ch {
		var got *io_prometheus_client.Metric
		m.Write(got)

		var want *io_prometheus_client.Metric
		prometheus.MustNewConstMetric(exporter.BoardIssueInfo, prometheus.GaugeValue, 1, "Test", "22", "358260205", "Backlog", "false").Write(want)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Collect() mismatch (-want +got):\n%s", diff)
		}
		close(ch)
	}
}

func ptrInt(n int) *int {
	return &n
}
