package dishwasher

import "context"

type Washer interface {
	Insert(ctx context.Context, id string, vector Vector, metadata Metadata) error
	Search(ctx context.Context, vector Vector) ([]Result, error)
}

type Vector []float32

type Metadata struct {
	Values map[string]any
}

type SearchOptions struct {
	TopK          int
	Ef            int
	Filter        map[string]any
	IncludeVector bool
}

type Result struct {
	ID       string
	Metadata map[string]any
	Vector   Vector
}
