package telepath

import "encoding/json"

func PackJSON(ctx *JSContext, value interface{}) (string, error) {
	newCtx := NewValueContext(ctx)
	v, err := newCtx.BuildNode(value)
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(
		v.Emit(),
	)
	return string(b), err
}
