---
description: Run comprehensive Go linting and static analysis tools
---

# Go Linting Suite

Run a comprehensive suite of Go linting and static analysis tools on the codebase.

**IMPORTANT:** This project requires `GOEXPERIMENT=jsonv2,synctest` to be set for all `go` commands. See `.claude/project-config.md` for details.

## Steps to Execute

1. **Run go vet**
   - Execute: `go vet ./...`
   - Report any issues with file paths and line numbers

2. **Run golangci-lint**
   - Execute: `go tool golangci-lint run`
   - Report any issues with file paths and line numbers
   - Highlight critical vs warning vs style issues

3. **Summarize results**
   - Count total issues by severity
   - Provide clickable file references using markdown format: `[filename.go:42](path/to/filename.go#L42)`
   - If all checks pass, clearly state "All linting checks passed ✓"
   - If there are issues, prioritize them by severity

## Output Format

Present results in this format:

```
## Linting Results

### go vet
[Status and issues]

### staticcheck
[Status and issues]

### golangci-lint
[Status and issues]

### Summary
[Total count and severity breakdown]
```

## Error Handling

- If a linter fails due to version compatibility, note this and continue
- Don't fail the entire check if one linter has issues
- Suggest fixes when appropriate (e.g., adding `//lint:ignore` comments)
