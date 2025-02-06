# `mf pull`

Pulls a workflow file from `modelflux.dev/hub` and saves it to the user's local machine to be run.

## Why

Allows users to quickly fetch any llm/agent workflow from a central `hub` so that they can
easily get started using and getting more familiar with AI-enabled workflows.

## Design

Uses a HTTP `GET` request with the workflow identifier to pull it from `modelflux.dev/hub`.

- Identifiers are in the format `user/workflowname:tag`
  - `user` is the user account that published the workflow
  - `workflowname` is the name of the workflow
  - `tag` is the version or specific tag that the workflow is pinned at. eg. `latest` or `1.2.3-beta` etc

The workflow is then saved to `$HOME/.modelflux/workflows/[user]/[workflowname]-[tag].yaml`.

Afterwards, the user can use `mf ls` to see this workflow and `mf run user/workflowname:tag` to run the workflow.
