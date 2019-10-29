package ctxx

import (
	"context"

	"github.com/fidelfly/gox/pkg/metax"
)

type metaKey struct{}

func WithMetadata(ctx context.Context, data ...map[interface{}]interface{}) context.Context {
	return context.WithValue(ctx, metaKey{}, metax.NewWrapMD(GetMetadata(ctx), data...))
}

func WrapMeta(ctx context.Context, md metax.MetaData) context.Context {
	return context.WithValue(ctx, metaKey{}, metax.Wrap(GetMetadata(ctx), md))
}

func GetMetadata(ctx context.Context) metax.MetaData {
	if v := ctx.Value(metaKey{}); v != nil {
		if md, ok := v.(metax.MetaData); ok {
			return md
		}
	}
	return nil
}
