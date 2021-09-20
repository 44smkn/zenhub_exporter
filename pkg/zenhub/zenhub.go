package zenhub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/44smkn/zenhub_exporter/pkg/model"
)

const (
	ZenhubRootEndpoint       = "https://api.zenhub.com"
	ZenhubWorkspacesEndpoint = "/p2/repositories/%s/workspaces"
	ZenhubBoardEndpoint      = "/p2/workspaces/%s/repositories/%s/board"
)

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

//go:generate mockgen -destination=../../mocks/zenhub/zenhub.go -package=mocks_zenhub github.com/44smkn/zenhub_exporter/pkg/zenhub Client
type Client interface {
	FetchWorkspaceIssues(ctx context.Context) ([]model.Issue, error)
}

func NewClient(token, repoID, workspaceName string, timeout time.Duration) Client {
	parsedURL, _ := url.ParseRequestURI(ZenhubRootEndpoint)
	httpClient := &http.Client{
		Timeout: timeout,
	}
	return &defaultClient{
		apiToken:      token,
		repoID:        repoID,
		workspaceName: workspaceName,
		URL:           parsedURL,
		httpClient:    httpClient,
	}
}

var _ Client = &defaultClient{}

type defaultClient struct {
	apiToken      string
	repoID        string
	workspaceName string

	URL        *url.URL
	httpClient *http.Client
}

func (c *defaultClient) newRequest(ctx context.Context, method, spath string, body io.Reader) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Authentication-Token", c.apiToken)
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (c *defaultClient) FetchWorkspaceIssues(ctx context.Context) ([]model.Issue, error) {
	workspace, err := c.fetchWorkspaces(ctx)
	if err != nil {
		return nil, err
	}

	issues := make([]model.Issue, 0, 200)

	for _, repoID := range workspace.Repositories {
		uri := fmt.Sprintf(ZenhubBoardEndpoint, workspace.ID, strconv.Itoa(repoID))
		req, err := c.newRequest(ctx, http.MethodGet, uri, nil)
		if err != nil {
			return nil, err
		}
		res, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		board := &ZenHubBoard{}
		if err := decodeBody(res, board); err != nil {
			return nil, err
		}

		for _, pipeline := range board.Pipelines {
			if pipeline.Issues == nil {
				continue
			}
			for _, is := range pipeline.Issues {
				var estimate *int
				if is.Estimate != nil {
					estimate = &is.Estimate.Value
				}
				issue := model.Issue{
					RepoID:        strconv.Itoa(repoID),
					WorkspaceName: c.workspaceName,
					IssueNumber:   strconv.Itoa(is.IssueNumber),
					Estimate:      estimate,
					IsEpic:        is.IsEpic,
					PipelineName:  pipeline.Name,
				}
				issues = append(issues, issue)
			}
		}
	}
	return issues, nil
}

func (c *defaultClient) fetchWorkspaces(ctx context.Context) (*ZenHubWorkspace, error) {
	uri := fmt.Sprintf(ZenhubWorkspacesEndpoint, c.repoID)
	req, err := c.newRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Requesting to %s is failed. Status: %s", uri, res.Status))
	}

	workspaces := make([]ZenHubWorkspace, 0, 10)
	if err := decodeBody(res, &workspaces); err != nil {
		return nil, err
	}

	for _, w := range workspaces {
		if w.Name == c.workspaceName {
			return &w, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Not contains specified workspace: %s", c.workspaceName))
}

// ZenHubBoard is the object representing board in zenhub world.
// It is bound with Get board data API Response of zenhub.
// See: https://github.com/ZenHubIO/API#get-a-zenhub-board-for-a-repository
type ZenHubBoard struct {
	Pipelines []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Issues []struct {
			IssueNumber int `json:"issue_number"`
			Estimate    *struct {
				Value int `json:"value"`
			} `json:"estimate,omitempty"`
			Position int  `json:"position,omitempty"`
			IsEpic   bool `json:"is_epic"`
		} `json:"issues"`
	} `json:"pipelines"`
}

// ZenHubWorkspace is the object representing workspace in zenhub world.
// It is bound with Get workspace data API Response of zenhub.
// See: https://github.com/ZenHubIO/API#get-zenhub-workspaces-for-a-repository
type ZenHubWorkspace struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ID           string `json:"id"`
	Repositories []int  `json:"repositories"`
}
