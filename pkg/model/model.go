package model

type Issue struct {
	RepoID        string
	WorkspaceName string
	IssueNumber   string
	Estimate      *int
	IsEpic        bool
	PipelineName  string
}
