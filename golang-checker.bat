
# 1️⃣ Run golangci-lint (Linting)
golangci-lint run ./...

# 2️⃣ Run go vet (Static Analysis)
go vet

# 4️⃣ Run gosec (Security Checks)
gosec ./...

# Measures cyclomatic complexity (helps identify complex functions).
gocyclo -over 10 ./...

# Official Go vulnerability scanner for dependencies.
govulncheck ./...


