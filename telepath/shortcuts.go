package telepath

import (
	"context"
	"encoding/json"
)

func PackJSON(ctx context.Context, context *JSContext, value interface{}) (string, error) {
	newCtx := NewValueContext(context)
	v, err := newCtx.BuildNode(ctx, value)
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(
		v.Emit(),
	)
	return string(b), err
}
