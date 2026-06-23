# SiYuan SDK Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build and release a Go SDK for the documented SiYuan local HTTP API.

**Architecture:** A single `Client` owns endpoint normalization, auth headers, transport, envelope decoding, and error semantics. Domain services hang off the client for typed APIs, while `Raw` remains available for newly added or unsupported endpoints.

**Tech Stack:** Go 1.22+ standard library only, `net/http`, `encoding/json`, `mime/multipart`, `httptest`, and `go test`.

## Global Constraints

- The Go module path is `github.com/lib-x/siyuan`.
- Public comments and code comments must be written in English.
- The SDK must avoid third-party dependencies for v0.1.0.
- Default endpoint is `http://127.0.0.1:6806`.
- API token auth uses `Authorization: Token <token>`.
- Service methods accept `context.Context` as the first argument.
- The release tag is `v0.1.0`.

---

### Task 1: Core Client, Options, and Errors

**Files:**
- Create: `go.mod`
- Create: `client.go`
- Create: `options.go`
- Create: `errors.go`
- Create: `raw.go`
- Test: `client_test.go`

**Interfaces:**
- Produces: `New(opts ...Option) (*Client, error)`, `WithEndpoint(string) Option`, `WithToken(string) Option`, `WithHTTPClient(*http.Client) Option`, `WithUserAgent(string) Option`, `WithHeader(string, string) Option`, `type APIError`.

- [x] **Step 1: Write tests for default client construction, endpoint normalization, auth headers, envelope unwrapping, API errors, and raw calls.**
- [x] **Step 2: Run `go test ./...` and confirm the new tests fail before implementation.**
- [x] **Step 3: Implement the minimal client, options, request helper, envelope decoder, and raw JSON call support.**
- [x] **Step 4: Run `go test ./...` and confirm the core tests pass.**

### Task 2: Typed Domain Services

**Files:**
- Create: `types.go`
- Create: `notebook.go`
- Create: `document.go`
- Create: `block.go`
- Create: `attribute.go`
- Create: `sql.go`
- Create: `template.go`
- Create: `export.go`
- Create: `convert.go`
- Create: `notification.go`
- Create: `network.go`
- Create: `system.go`
- Test: `services_test.go`

**Interfaces:**
- Consumes: `(*Client).postJSON(ctx context.Context, path string, input any, output any) error`.
- Produces: typed services available from `Client`, including `Notebooks`, `Documents`, `Blocks`, `Attributes`, `SQL`, `Templates`, `Export`, `Convert`, `Notifications`, `Network`, and `System`.

- [x] **Step 1: Write service tests through `httptest.Server` for representative no-arg, JSON request, error, and typed response flows.**
- [x] **Step 2: Implement request and response types matching the official API fields.**
- [x] **Step 3: Implement all documented JSON POST service methods.**
- [x] **Step 4: Run `go test ./...` and confirm service coverage passes.**

### Task 3: Multipart and File APIs

**Files:**
- Create: `asset.go`
- Create: `file.go`
- Test: `multipart_file_test.go`

**Interfaces:**
- Consumes: `(*Client).newRequest(ctx context.Context, method string, path string, body io.Reader) (*http.Request, error)`.
- Produces: `Assets.Upload`, `Files.Get`, `Files.Put`, `Files.Remove`, `Files.Rename`, and `Files.ReadDir`.

- [x] **Step 1: Write multipart tests that inspect form fields, uploaded filenames, and auth headers.**
- [x] **Step 2: Implement multipart upload helpers using `io.Reader` inputs.**
- [x] **Step 3: Implement `/api/file/getFile` status handling for `200` file content and `202` JSON errors.**
- [x] **Step 4: Run `go test ./...` and confirm multipart and file tests pass.**

### Task 4: README, Examples, and Release

**Files:**
- Create: `README.md`
- Create: `examples/basic/main.go`
- Create: `examples/assets/main.go`
- Create: `examples/raw/main.go`

**Interfaces:**
- Consumes: public SDK APIs from Tasks 1-3.
- Produces: copy-pasteable usage documentation and runnable example programs.

- [x] **Step 1: Write README with install command, quick start, auth, endpoint configuration, service map, error handling, upload examples, and release notes.**
- [x] **Step 2: Add examples that compile with the module path `github.com/lib-x/siyuan`.**
- [x] **Step 3: Run `gofmt`, `go test ./...`, and `go test -race ./...`.**
- [ ] **Step 4: Commit the SDK and tag `v0.1.0`.**
