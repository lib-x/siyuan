# SiYuan Go SDK

Go SDK for the SiYuan local HTTP API.

The SDK wraps SiYuan's JSON envelope, token authentication, multipart uploads, file APIs, and typed service groups while keeping a `Raw` escape hatch for endpoints that are not yet covered.

## Install

```bash
go get github.com/lib-x/siyuan@v0.1.0
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lib-x/siyuan"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := siyuan.New(
		siyuan.WithToken(os.Getenv("SIYUAN_TOKEN")),
	)
	if err != nil {
		log.Fatal(err)
	}

	version, err := client.System.Version(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SiYuan:", version)

	notebooks, err := client.Notebooks.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, notebook := range notebooks {
		fmt.Printf("%s %s\n", notebook.ID, notebook.Name)
	}
}
```

By default the client uses `http://127.0.0.1:6806`.

```go
client, err := siyuan.New(
	siyuan.WithEndpoint("http://127.0.0.1:6806"),
	siyuan.WithToken(os.Getenv("SIYUAN_TOKEN")),
)
```

## Authentication

Find the API token in SiYuan under **Settings - About** and pass it with `WithToken`.

The SDK sends:

```text
Authorization: Token <token>
```

## Services

| Service | Field | Examples |
|---|---|---|
| Notebooks | `client.Notebooks` | `List`, `Create`, `Open`, `Close`, `Rename`, `Remove`, `GetConf`, `SetConf` |
| Documents | `client.Documents` | `CreateWithMarkdown`, `Rename`, `Remove`, `Move`, `GetHPathByID`, `GetPathByID` |
| Assets | `client.Assets` | `Upload` |
| Blocks | `client.Blocks` | `Insert`, `Append`, `Update`, `Delete`, `Move`, `GetKramdown`, `GetChildBlocks` |
| Attributes | `client.Attributes` | `SetBlockAttrs`, `GetBlockAttrs` |
| SQL | `client.SQL` | `Query`, `FlushTransaction` |
| Templates | `client.Templates` | `Render`, `RenderSprig` |
| Files | `client.Files` | `Get`, `Put`, `Remove`, `Rename`, `ReadDir` |
| Export | `client.Export` | `MarkdownContent`, `Resources` |
| Convert | `client.Convert` | `Pandoc` |
| Notifications | `client.Notifications` | `PushMsg`, `PushErrMsg` |
| Network | `client.Network` | `ForwardProxy` |
| System | `client.System` | `BootProgress`, `Version`, `CurrentTime` |
| Raw | `client.Raw` | `Post` |

## Create a Document

```go
docID, err := client.Documents.CreateWithMarkdown(ctx, siyuan.CreateDocWithMarkdownRequest{
	Notebook: notebookID,
	Path:     "/api demo/hello",
	Markdown: "# Hello from Go\n\nCreated with github.com/lib-x/siyuan.",
})
if err != nil {
	return err
}
fmt.Println(docID)
```

Calling SiYuan with the same human-readable path does not overwrite an existing document.

## Insert a Block

```go
ops, err := client.Blocks.Insert(ctx, siyuan.InsertBlockRequest{
	DataType:   siyuan.DataTypeMarkdown,
	Data:       "A new paragraph",
	PreviousID: previousBlockID,
})
if err != nil {
	return err
}
fmt.Println(ops[0].DoOperations[0].ID)
```

`NextID`, `PreviousID`, or `ParentID` must identify where SiYuan should insert the block. SiYuan prioritizes `NextID`, then `PreviousID`, then `ParentID`.

## Upload Assets

```go
file, err := os.Open("image.png")
if err != nil {
	return err
}
defer file.Close()

result, err := client.Assets.Upload(ctx, siyuan.UploadAssetsRequest{
	AssetsDirPath: "/assets/",
	Files: []siyuan.UploadAssetFile{
		{Name: "image.png", Reader: file},
	},
})
if err != nil {
	return err
}
fmt.Println(result.SuccMap["image.png"])
```

## File APIs

`Files.Get` is special because SiYuan returns raw file content with HTTP `200`, and a JSON error envelope with HTTP `202`.

```go
content, err := client.Files.Get(ctx, "/data/assets/example.txt")
if err != nil {
	return err
}
fmt.Println(string(content))
```

## Raw API

Use `Raw.Post` when SiYuan adds a new endpoint before the SDK exposes a typed method.

```go
var currentTime int64
err := client.Raw.Post(ctx, "/api/system/currentTime", nil, &currentTime)
```

`Raw.Post` still unwraps SiYuan's response envelope and returns only `data`.

## Error Handling

SiYuan JSON endpoints return:

```json
{
  "code": 0,
  "msg": "",
  "data": {}
}
```

The SDK returns the unwrapped `data`. If `code` is non-zero, methods return `*siyuan.APIError`.

```go
var apiErr *siyuan.APIError
if errors.As(err, &apiErr) {
	fmt.Printf("SiYuan error: code=%d msg=%s path=%s\n", apiErr.Code, apiErr.Message, apiErr.Path)
}
```

Unexpected HTTP statuses return `*siyuan.HTTPError`.

## Configuration

```go
client, err := siyuan.New(
	siyuan.WithEndpoint("http://127.0.0.1:6806"),
	siyuan.WithToken(token),
	siyuan.WithHTTPClient(&http.Client{Timeout: 15 * time.Second}),
	siyuan.WithUserAgent("my-tool/1.0"),
	siyuan.WithHeader("X-Trace-ID", traceID),
)
```

The SDK does not retry requests by default. Many SiYuan APIs mutate data, so automatic retries can create duplicate documents, duplicate blocks, or repeated deletes. Add retry behavior through a custom `http.Client` only when the called operation is safe for your workflow.

## Examples

```bash
SIYUAN_TOKEN=xxx go run ./examples/basic
SIYUAN_TOKEN=xxx go run ./examples/raw
SIYUAN_TOKEN=xxx go run ./examples/assets ./image.png
```

Set `SIYUAN_ENDPOINT` to override `http://127.0.0.1:6806`.

## Development

```bash
go test ./...
go test -race ./...
```

## Release v0.1.0

Initial release:

- Typed client with functional options.
- Uniform SiYuan envelope decoding.
- `APIError` and `HTTPError`.
- Typed services for the documented official API.
- Multipart asset upload and file write support.
- Raw fallback for unsupported endpoints.
