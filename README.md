# zenhub_exporter

This is a simple server that scrapes [ZenHub](https://www.zenhub.com/) stats and exports them via HTTP for Prometheus consumption.

## Installation

The `zenhub_exporter` listens on HTTP port 9861 by default. See the `--help` output for more options.

You must pass below items on to zenhub_exporter:

* pass zenhub api key via `--zenhub.api-token` command flag or `ZENHUB_API_TOKEN` environemnt variables
* pass zenhub workspace name via `--zenhub.workspace-id` command flag or `ZENHUB_WORKSPACE_NAME` environment variables
* pass zenhub repository id bound with board via `--zenhub.repository-id` command flag or `ZENHUB_REPO_ID` environment variables

Executing below command, your `repository-id` will be found.

```console
# replace org and repo with yours
$ gh api /repos/{org}/{repo} --jq '.id'

# for example
$ gh api /repos/44smkn/zenhub_exporter --jq '.id'
404765867
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: zenhub-exporter
  labels:
    app: zenhub-exporter
spec:
  replicas: 1
  selector:
    matchLabels:
      app: zenhub-exporter
  template:
    metadata:
      labels:
        app: zenhub-exporter
    spec:
      containers:
      - name: zenhub-exporter
        image: ghcr.io/44smkn/zenhub_exporter:latest
        ports:
        - containerPort: 9861
        env:
        - name: ZENHUB_API_TOKEN
          valueFrom:
            secretKeyRef:
              name: zenhub-secret
              key: api-token
        - name: ZENHUB_WORKSPACE_NAME
          value: Test
        - name: ZEBNHUB_REPO_ID
          value: 404765867
---
apiVersion: v1
kind: Secret
metadata:
  name: zenhub-secret
type: Opaque
data:
  api-token: YWRtaW4=
```

## Collectors

| Metric Name                                   | Metric type | Description                              | Labels |
| --------------------------------------------- | ----------- | ---------------------------------------- | ------ |
| zenhub_workspace_issue_info                   | Gauge       | Information about issue managed by board | `workspace`=&lt;workspace-name&gt;<br>`issue_number`=&lt;issue-number&gt;<br>`repository_id`=&lt;repository-id&gt;<br>`pipeline`=&lt;pipeline-name&gt;<br>`is_epic`=&lt;is_epic&gt;|
| zenhub_workspace_issue_estimated_story_points | Gauge       | Estimated story point of each issue      | `workspace`=&lt;workspace-name&gt;<br>`issue_number`=&lt;issue-number&gt;<br>`repository_id`=&lt;repository-id&gt;|
