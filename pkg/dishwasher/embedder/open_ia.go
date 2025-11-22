package embedder

import (
	"context"
	"fmt"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/tomiok/dishwasher/pkg/dishwasher"
	"time"
)

type OpenAI struct {
	APIKey    string
	Model     string
	Dimension int
	client    openai.Client
}

func NewOpenAIEmbedder(apiKey, model string) (OpenAI, error) {
	if apiKey == "" {
		return OpenAI{}, fmt.Errorf("APIKey required")
	}
	if model == "" {
		model = "text-embedding-3-small" // sensible default
	}

	client := openai.NewClient(
		option.WithRequestTimeout(time.Second),
		option.WithAPIKey(apiKey),
	)

	// Determine dimension based on model
	var dimension int
	switch model {
	case "text-embedding-3-small":
		dimension = 1536
	case "text-embedding-3-large":
		dimension = 3072
	default:
		return OpenAI{}, fmt.Errorf("unknown model: %s", model)
	}

	return OpenAI{
		APIKey:    apiKey,
		Model:     model,
		client:    client,
		Dimension: dimension,
	}, nil
}

func (o OpenAI) Embed(ctx context.Context, text string) (dishwasher.Vector64, error) {
	res, err := o.EmbedBatch(ctx, []string{text})
	if err != nil {
		return nil, err
	}

	return res[0], nil
}

func (o OpenAI) EmbedBatch(ctx context.Context, texts []string) ([]dishwasher.Vector64, error) {
	response, err := o.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: texts,
		},
		Model:      o.Model,
		Dimensions: openai.Int(int64(o.Dimension)),
	})

	if err != nil {
		return nil, fmt.Errorf("openai embed failed: %w", err)
	}

	embeddings := make([]dishwasher.Vector64, len(response.Data))
	for i, data := range response.Data {
		embeddings[i] = data.Embedding
	}

	return embeddings, nil
}

func (o OpenAI) Dimensions() int {
	//TODO implement me
	panic("implement me")
}

func (o OpenAI) Name() string {
	//TODO implement me
	panic("implement me")
}
