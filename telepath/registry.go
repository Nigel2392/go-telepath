package telepath

import (
	"reflect"
)

var (
	specificAdapterMap = make(map[reflect.Kind]map[reflect.Type]Adapter)
	defaultAdapterMap  = make(map[reflect.Kind]Adapter)
	iFaceAdapterMap    = make(map[reflect.Type]Adapter)
)

func init() {
	var (
		rTypBool    = reflect.TypeOf(bool(false))
		rTypInt     = reflect.TypeOf(int(0))
		rTypInt8    = reflect.TypeOf(int8(0))
		rTypInt16   = reflect.TypeOf(int16(0))
		rTypInt32   = reflect.TypeOf(int32(0))
		rTypInt64   = reflect.TypeOf(int64(0))
		rTypUint    = reflect.TypeOf(uint(0))
		rTypUint8   = reflect.TypeOf(uint8(0))
		rTypUint16  = reflect.TypeOf(uint16(0))
		rTypUint32  = reflect.TypeOf(uint32(0))
		rTypUint64  = reflect.TypeOf(uint64(0))
		rTypFloat32 = reflect.TypeOf(float32(0))
		rTypFloat64 = reflect.TypeOf(float64(0))
		rTypString  = reflect.TypeOf(string(""))
		rTypSlice   = reflect.TypeOf([]interface{}{})
		rTypMap     = reflect.TypeOf(map[string]interface{}{})

		// Interface types
		rTypError = reflect.TypeOf((*error)(nil)).Elem()
	)

	specificAdapterMap[rTypBool.Kind()] = make(map[reflect.Type]Adapter)
	specificAdapterMap[rTypInt.Kind()] = make(map[reflect.Type]Adapter)
	specificAdapterMap[rTypInt8.Kind()] = make(map[reflect.Type]Adapter)
	specificAdapterMap[rTypInt16.Kind()] = make(map[reflect.Type]Adapter)
	specificAdapterMap[rTypInt32.Kind()] = make(map[reflect.Type]Adapter)
	specificAdapterMap[rTypInt64.Kind()] = make(map[reflect.Type]Adapter)
	specificAdapterMap[rTypUint.Kind()] = make(map[reflect.Type]Adapter)
	specificAdapterMap[rTypUint8.Kind()] = make(map[reflect.Type]Adapter)
	specificAdapterMap[rTypUint16.Kind()] = make(map[reflect.Type]Adapter)
	specificAdapterMap[rTypUint32.Kind()] = make(map[reflect.Type]Adapter)
	specificAdapterMap[rTypUint64.Kind()] = make(map[reflect.Type]Adapter)
	specificAdapterMap[rTypFloat32.Kind()] = make(map[reflect.Type]Adapter)
	specificAdapterMap[rTypFloat64.Kind()] = make(map[reflect.Type]Adapter)
	specificAdapterMap[rTypString.Kind()] = make(map[reflect.Type]Adapter)
	specificAdapterMap[rTypSlice.Kind()] = make(map[reflect.Type]Adapter)
	specificAdapterMap[rTypMap.Kind()] = make(map[reflect.Type]Adapter)

	specificAdapterMap[rTypBool.Kind()][rTypBool] = BaseAdapter()
	specificAdapterMap[rTypInt.Kind()][rTypInt] = BaseAdapter()
	specificAdapterMap[rTypInt8.Kind()][rTypInt8] = BaseAdapter()
	specificAdapterMap[rTypInt16.Kind()][rTypInt16] = BaseAdapter()
	specificAdapterMap[rTypInt32.Kind()][rTypInt32] = BaseAdapter()
	specificAdapterMap[rTypInt64.Kind()][rTypInt64] = BaseAdapter()
	specificAdapterMap[rTypUint.Kind()][rTypUint] = BaseAdapter()
	specificAdapterMap[rTypUint8.Kind()][rTypUint8] = BaseAdapter()
	specificAdapterMap[rTypUint16.Kind()][rTypUint16] = BaseAdapter()
	specificAdapterMap[rTypUint32.Kind()][rTypUint32] = BaseAdapter()
	specificAdapterMap[rTypUint64.Kind()][rTypUint64] = BaseAdapter()
	specificAdapterMap[rTypFloat32.Kind()][rTypFloat32] = BaseAdapter()
	specificAdapterMap[rTypFloat64.Kind()][rTypFloat64] = BaseAdapter()
	specificAdapterMap[rTypString.Kind()][rTypString] = StringAdapter()
	specificAdapterMap[rTypSlice.Kind()][rTypSlice] = SliceAdapter()
	specificAdapterMap[rTypMap.Kind()][rTypMap] = MapAdapter()

	defaultAdapterMap[rTypBool.Kind()] = BaseAdapter()
	defaultAdapterMap[rTypInt.Kind()] = BaseAdapter()
	defaultAdapterMap[rTypInt8.Kind()] = BaseAdapter()
	defaultAdapterMap[rTypInt16.Kind()] = BaseAdapter()
	defaultAdapterMap[rTypInt32.Kind()] = BaseAdapter()
	defaultAdapterMap[rTypInt64.Kind()] = BaseAdapter()
	defaultAdapterMap[rTypUint.Kind()] = BaseAdapter()
	defaultAdapterMap[rTypUint8.Kind()] = BaseAdapter()
	defaultAdapterMap[rTypUint16.Kind()] = BaseAdapter()
	defaultAdapterMap[rTypUint32.Kind()] = BaseAdapter()
	defaultAdapterMap[rTypUint64.Kind()] = BaseAdapter()
	defaultAdapterMap[rTypFloat32.Kind()] = BaseAdapter()
	defaultAdapterMap[rTypFloat64.Kind()] = BaseAdapter()
	defaultAdapterMap[rTypString.Kind()] = StringAdapter()
	defaultAdapterMap[rTypSlice.Kind()] = SliceAdapter()
	defaultAdapterMap[rTypMap.Kind()] = MapAdapter()

	// Interface adapters
	iFaceAdapterMap[rTypError] = ErrorAdapter()
}

type AdapterRegistry struct {
	adapters map[reflect.Kind]map[reflect.Type]Adapter
	defaults map[reflect.Kind]Adapter
	iFaces   map[reflect.Type]Adapter
}

func NewAdapterRegistry() *AdapterRegistry {
	return &AdapterRegistry{
		adapters: specificAdapterMap,
		defaults: defaultAdapterMap,
		iFaces:   iFaceAdapterMap,
	}
}

func (r *AdapterRegistry) RegisterAdapter(k reflect.Kind, t reflect.Type, a Adapter) {
	if _, ok := r.adapters[k]; !ok {
		r.adapters[k] = make(map[reflect.Type]Adapter)
	}

	r.adapters[k][t] = a
}

func (r *AdapterRegistry) RegisterDefaultAdapter(k reflect.Kind, a Adapter) {
	r.defaults[k] = a
}

func (r *AdapterRegistry) Context() *JSContext {
	var c = &JSContext{
		Media:           &nullMedia{},
		AdapterRegistry: r,
	}
	return c
}

func (r *AdapterRegistry) RegisterInterfaceAdapter(a Adapter, i interface{}) {
	var t = reflect.TypeOf(i)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Interface {
		panic("RegisterInterfaceAdapter: i must be an interface")
	}

	r.iFaces[t] = a
}

func (r *AdapterRegistry) Register(adapter any, forType ...interface{}) {
	var v interface{}

	if len(forType) == 0 {
		v = adapter
	} else if len(forType) == 1 {
		v = forType[0]
	} else {
		for _, t := range forType {
			r.Register(adapter, t)
		}
		return
	}

	if getter, ok := v.(AdapterGetter); ok {
		adapter = getter.Adapter()
	}

	var (
		a = adapter.(Adapter)
		t = reflect.TypeOf(v)
		k = t.Kind()
	)

	r.RegisterAdapter(k, t, a)

	// If the type is a pointer, register the adapter for the underlying type as well
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		r.RegisterAdapter(t.Kind(), t, a)
	}
}

func (r *AdapterRegistry) Find(value interface{}) (Adapter, bool) {

	if value == nil {
		return nil, false
	}

	if getter, ok := value.(AdapterGetter); ok {
		var a = getter.Adapter()
		if a != nil {
			return a, true
		}
	}

	var (
		v = reflect.ValueOf(value)
		k = v.Kind()
	)

	if k == reflect.Invalid {
		return nil, false
	}

	var (
		t     = v.Type()
		iType reflect.Type
		a     Adapter
		ok    bool
	)

	if _, ok = r.adapters[k]; ok {
		if a, ok = r.adapters[k][t]; ok {
			return a, true
		}
	}

	for iType, a = range r.iFaces {
		if t.Implements(iType) {
			return a, true
		}
	}

	if a, ok = r.defaults[k]; ok {
		return a, true
	}

	return nil, false
}
