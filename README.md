# Dash Prototype

To run, first log into your kubernetes cluster, and:

```
export GO111MODULE=on
go mod vendor
# Run a dash inventory
go run cmd/dash.go -i examples/v3/
# Run an existing applier inventory
go run cmd/dash.go -i examples/v2/
```
