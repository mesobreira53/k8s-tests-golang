
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run ./...



go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...


go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
gocyclo -over 10 ./...

go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

