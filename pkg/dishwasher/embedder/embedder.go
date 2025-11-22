package embedder

import (
	"context"
	"github.com/tomiok/dishwasher/pkg/dishwasher"
)

type Embedder interface {
	Embed(ctx context.Context, text string) (dishwasher.Vector, error)
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
	Dimensions() int
	Name() string
}
