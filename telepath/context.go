package telepath

import (
	"fmt"
	"reflect"
)

var _ Context = (*ValueContext)(nil)

type JSContext struct {
	Media           Media
	AdapterRegistry *AdapterRegistry
}

func (c *JSContext) AddMedia(media Media) {
	c.Media = media.Merge(c.Media)
}

func (c *JSContext) Registry() *AdapterRegistry {
	return c.AdapterRegistry
}

func (c *JSContext) Pack(value interface{}) (interface{}, error) {
	var newCtx = NewValueContext(c)
	var v, err = newCtx.BuildNode(value)
	if err != nil {
		return nil, err
	}
	return v.Emit(), nil
}

type ValueContext struct {
	ParentContext   *JSContext
	AdapterRegistry *AdapterRegistry
	Nodes           map[uintptr]Node
	RawValues       map[uintptr]interface{} // keep reference to prevent GC
	NextID          int
}

func NewValueContext(c *JSContext) *ValueContext {
	return &ValueContext{
		ParentContext:   c,
		AdapterRegistry: c.Registry(),
		Nodes:           make(map[uintptr]Node),
		RawValues:       make(map[uintptr]interface{}),
	}
}

func (c *ValueContext) AddMedia(media Media) {
	c.ParentContext.AddMedia(media)
}

func (c *ValueContext) Registry() *AdapterRegistry {
	return c.AdapterRegistry
}

func (c *ValueContext) buildNewNode(value interface{}) (Node, error) {

	var adapter, ok = c.AdapterRegistry.Find(value)
	if ok {
		return adapter.BuildNode(value, c)
	}

	var v = reflect.ValueOf(value)

	if !v.IsValid() {
		return NullNode(), nil
	}

	switch v.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool, reflect.Invalid:
		return BaseAdapter().BuildNode(value, c)
	case reflect.Slice, reflect.Array:
		var n, err = SliceAdapter().BuildNode(value, c)
		return n, err
	case reflect.Map:
		return MapAdapter().BuildNode(value, c)
	case reflect.String:
		return StringAdapter().BuildNode(value, c)
	}

	return nil, fmt.Errorf("no adapter found for value %v (%T)", value, value)
}

func (c *ValueContext) BuildNode(value interface{}) (Node, error) {
	var (
		rVal   = reflect.ValueOf(value)
		node   Node
		objKey uintptr
		ok     bool
	)

	switch rVal.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map:
		objKey = rVal.Pointer()
	}

	if node, ok = c.Nodes[objKey]; !ok {
		node, err := c.buildNewNode(value)
		if err != nil {
			return nil, err
		}
		if objKey != 0 {
			c.Nodes[objKey] = node
			c.RawValues[objKey] = value
		}
		return node, nil
	}

	if node.GetID() == 0 {
		c.NextID++
		node.SetID(c.NextID)
	}

	return node, nil
}
