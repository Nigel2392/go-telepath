package telepath

import (
	"context"
	"html/template"

	"golang.org/x/exp/constraints"
)

var DICT_RESERVED_KEYS = []string{
	"_args",
	"_dict",
	"_id",
	"_list",
	"_ref",
	"_type",
	"_val",
}

const STRING_REF_MIN_LENGTH = 20 // Strings shorter than this will not be turned into references

type TelepathValue struct {
	Type string                 `json:"_type,omitempty"`
	Args []any                  `json:"_args,omitempty"`
	Dict map[string]interface{} `json:"_dict,omitempty"`
	List []interface{}          `json:"_list,omitempty"`
	Val  interface{}            `json:"_val,omitempty"`
	Ref  int                    `json:"_ref,omitempty"`
	ID   int                    `json:"_id,omitempty"`
}

type AdapterGetter interface {
	Adapter(ctx context.Context) Adapter
}

type Adapter interface {
	BuildNode(ctx context.Context, value interface{}, context Context) (Node, error)
}

type Context interface {
	AddMedia(media Media)
	BuildNode(ctx context.Context, value interface{}) (Node, error)
	Registry() *AdapterRegistry
}

type Media interface {
	Merge(other Media) Media
	CSS() []template.HTML
	JS() []template.HTML
}

// If this node is assigned an id, emit() should return the verbose representation with the
// id attached on first call, and a reference on subsequent calls. To disable this behaviour
// (e.g. for small primitive values where the reference representation adds unwanted overhead),
// set self.use_id = False.
type Node interface {
	Emit() any                  // emit (returns a dict representation of a value, this should be the main method used by an application.)
	EmitVerbose() TelepathValue // emit_verbose (returns a dict representation of a value that can have an _id attached)
	EmitCompact() any           // emit_compact (returns a compact representation of the value, in any JSON-serialisable type)
	GetValue() interface{}

	UseID() bool
	SetID(id int)
	GetID() int
}

type PrimitiveNodeValue interface {
	constraints.Integer | constraints.Float | bool
}

var (
	GlobalRegistry    = NewAdapterRegistry()
	NewContext        = GlobalRegistry.Context
	Register          = GlobalRegistry.Register
	RegisterInterface = GlobalRegistry.RegisterInterfaceAdapter
	Find              = GlobalRegistry.Find
)
