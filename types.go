package siyuan

// NotebookID identifies a SiYuan notebook.
type NotebookID string

// DocumentID identifies a SiYuan document.
type DocumentID string

// BlockID identifies a SiYuan block.
type BlockID string

// DataType is the payload format used by block mutation APIs.
type DataType string

const (
	// DataTypeMarkdown sends Markdown content.
	DataTypeMarkdown DataType = "markdown"
	// DataTypeDOM sends SiYuan DOM content.
	DataTypeDOM DataType = "dom"
)

// OperationTransaction describes operations returned by block mutation APIs.
type OperationTransaction struct {
	DoOperations   []Operation `json:"doOperations"`
	UndoOperations []Operation `json:"undoOperations"`
}

// Operation describes one block operation returned by SiYuan.
type Operation struct {
	Action     string    `json:"action"`
	Data       any       `json:"data"`
	ID         BlockID   `json:"id"`
	ParentID   BlockID   `json:"parentID"`
	PreviousID BlockID   `json:"previousID"`
	NextID     BlockID   `json:"nextID"`
	RetData    any       `json:"retData"`
	SrcIDs     []BlockID `json:"srcIDs"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
}
