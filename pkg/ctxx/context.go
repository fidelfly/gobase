package ctxx

import (
	"context"

	"github.com/fidelfly/gox/pkg/metax"
)

type metaKey struct{}

func WithMetadata(ctx context.Context, data ...map[interface{}]interface{}) context.Context {
	md := GetMetadata(ctx)
	if md != nil {
		metax.FillMd(md, data...)
		return ctx
	}
	return context.WithValue(ctx, metaKey{}, metax.NewMD(data...))
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

func IsCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func IsCancellable(ctx context.Context) (<-chan struct{}, bool) {
	dch := ctx.Done()
	return dch, dch != nil
}
