# tbd

Recommend getting the go extension for vscode.
Go has a built in linter/formatter

Get go and make if you don't already have them

```bash
brew install go make
```

Run the application with

```bash
go run ./cmd
# Eg. use the 'pull' command
go run ./cmd pull
```

Build the application with

```bash
# in the project directory
make
# run binary directly
./_output/bin/tbd
```

## Project Plan

### CLI Commands

- pull
- run
- list, ls
- push (future)
- explore (future)

### Nodes (each with inputs/outputs)

- Base models
  - Azure
  - OpenAI
  - Local (Ollama serve)
- Memory store ?
- Model actions
  - Summary
  - Decision
  - Monitor (agent)
- Tools
- Triggers
  - Time based
  - API based?
- Conditionals
- Loops

### Sample workflow/compose .yaml file

```yaml
models:
  model-1:
    type: openai
    key: shxxxxxxx
    model-class: 4o
  model-2:
    type: azure-openai
    key: xdkfjksdjf
    endpoint: http://...../v1/
tools:
  tool-1:
    type: web-browser-v0.2
trigger:
  time:
    at: 2025-03-21 22:33:33
    repeat: daily
workflow:
```

### GUI

React flow based node editing
