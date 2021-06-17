package chunck

import (
	"github.com/ITLab-Projects/pkg/chunkresp"
	"context"
	"errors"
)

var (
	ErrChunckNotFound = errors.New("Chunck not found in context")
)

type ChunkKey struct{}

type ChunckContext struct {
	context.Context
}

func New(ctx context.Context, chunck *chunkresp.ChunkResp) *ChunckContext {
	return &ChunckContext{
		Context: context.WithValue(
			ctx,
			ChunkKey{},
			chunck,
		),
	}
}

func GetChunckFromContext(ctx context.Context) (*chunkresp.ChunkResp, error) {
	val := ctx.Value(ChunkKey{})
	if chunck, ok := val.(*chunkresp.ChunkResp); !ok {
		return nil, ErrChunckNotFound
	} else {
		return chunck, nil
	}
}