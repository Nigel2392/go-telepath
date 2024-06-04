package telepath

import "encoding/json"

func PackJSON(ctx *JSContext, value interface{}) ([]byte, error) {
	var newCtx = NewValueContext(ctx)
	var v, err = newCtx.BuildNode(value)
	if err != nil {
		return nil, err
	}
	return json.Marshal(
		v.Emit(),
	)
}
