package siyuan

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotebookServiceList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/notebook/lsNotebooks" {
			t.Fatalf("path = %s, want /api/notebook/lsNotebooks", r.URL.Path)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code": 0,
			"msg":  "",
			"data": map[string]any{
				"notebooks": []map[string]any{
					{
						"id":     "20210817205410-2kvfpfn",
						"name":   "test",
						"icon":   "1f4d4",
						"sort":   1,
						"closed": false,
					},
				},
			},
		})
	}))
	defer server.Close()

	client, err := New(WithEndpoint(server.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	notebooks, err := client.Notebooks.List(context.Background())
	if err != nil {
		t.Fatalf("Notebooks.List() error = %v", err)
	}
	if len(notebooks) != 1 {
		t.Fatalf("len(notebooks) = %d, want 1", len(notebooks))
	}
	if notebooks[0].ID != NotebookID("20210817205410-2kvfpfn") || notebooks[0].Name != "test" {
		t.Fatalf("notebook = %#v, want id and name decoded", notebooks[0])
	}
}

func TestDocumentServiceCreateWithMarkdown(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/filetree/createDocWithMd" {
			t.Fatalf("path = %s, want /api/filetree/createDocWithMd", r.URL.Path)
		}

		var body CreateDocWithMarkdownRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Decode request body: %v", err)
		}
		if body.Notebook != NotebookID("box") || body.Path != "/foo/bar" || body.Markdown != "# Title" {
			t.Fatalf("request body = %#v, want create doc payload", body)
		}

		writeJSON(t, w, http.StatusOK, map[string]any{
			"code": 0,
			"msg":  "",
			"data": "20210914223645-oj2vnx2",
		})
	}))
	defer server.Close()

	client, err := New(WithEndpoint(server.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	id, err := client.Documents.CreateWithMarkdown(context.Background(), CreateDocWithMarkdownRequest{
		Notebook: NotebookID("box"),
		Path:     "/foo/bar",
		Markdown: "# Title",
	})
	if err != nil {
		t.Fatalf("Documents.CreateWithMarkdown() error = %v", err)
	}
	if id != DocumentID("20210914223645-oj2vnx2") {
		t.Fatalf("id = %q, want created document id", id)
	}
}

func TestBlockServiceInsert(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/block/insertBlock" {
			t.Fatalf("path = %s, want /api/block/insertBlock", r.URL.Path)
		}

		var body InsertBlockRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Decode request body: %v", err)
		}
		if body.DataType != DataTypeMarkdown || body.Data != "hello" || body.PreviousID != BlockID("prev") {
			t.Fatalf("request body = %#v, want insert payload", body)
		}

		writeJSON(t, w, http.StatusOK, map[string]any{
			"code": 0,
			"msg":  "",
			"data": []map[string]any{
				{
					"doOperations": []map[string]any{
						{
							"action":     "insert",
							"id":         "new-block",
							"data":       "<div></div>",
							"parentID":   "",
							"previousID": "prev",
						},
					},
					"undoOperations": nil,
				},
			},
		})
	}))
	defer server.Close()

	client, err := New(WithEndpoint(server.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	transactions, err := client.Blocks.Insert(context.Background(), InsertBlockRequest{
		DataType:   DataTypeMarkdown,
		Data:       "hello",
		PreviousID: BlockID("prev"),
	})
	if err != nil {
		t.Fatalf("Blocks.Insert() error = %v", err)
	}
	if len(transactions) != 1 || len(transactions[0].DoOperations) != 1 {
		t.Fatalf("transactions = %#v, want one insert operation", transactions)
	}
	if transactions[0].DoOperations[0].ID != BlockID("new-block") {
		t.Fatalf("operation id = %q, want new-block", transactions[0].DoOperations[0].ID)
	}
}

func TestSystemServiceVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/system/version" {
			t.Fatalf("path = %s, want /api/system/version", r.URL.Path)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code": 0,
			"msg":  "",
			"data": "3.0.0",
		})
	}))
	defer server.Close()

	client, err := New(WithEndpoint(server.URL))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	version, err := client.System.Version(context.Background())
	if err != nil {
		t.Fatalf("System.Version() error = %v", err)
	}
	if version != "3.0.0" {
		t.Fatalf("version = %q, want 3.0.0", version)
	}
}
