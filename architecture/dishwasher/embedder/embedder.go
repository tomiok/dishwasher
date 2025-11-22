package embedder

import "dishwasher/architecture/dishwasher"

type Embedder interface {
	Embed(text string) (dishwasher.Vector, error)
}
