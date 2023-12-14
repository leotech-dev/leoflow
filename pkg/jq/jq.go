package jq

import (
	"context"

	"github.com/itchyny/gojq"
)

type Jq func(ctx context.Context, v any, values ...any) gojq.Iter

func New(expr string) (Jq, error) {
	query, err := gojq.Parse(expr)
	if err != nil {
		return nil, err
	}

	code, err := gojq.Compile(query, gojq.WithFunction("fetch", 1, 2, fetch))
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context, v any, values ...any) gojq.Iter {
		return code.RunWithContext(ctx, v, values...)
	}, nil
}
