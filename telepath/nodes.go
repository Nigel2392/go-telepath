package telepath

import (
	"reflect"
	"slices"

	"github.com/google/uuid"
)

var _ Node = (*TelepathNode)(nil)

type TelepathNode struct {
	ID            int
	Seen          bool
	UseIdentifier bool
	EmitVerboseFn func() TelepathValue
	EmitCompactFn func() any
}

func NewTelepathNode() *TelepathNode {
	return &TelepathNode{}
}

func (m *TelepathNode) GetValue() interface{} {
	return nil
}

func (m *TelepathNode) SetID(id int) {
	m.ID = id
	m.UseIdentifier = true
}

func (m *TelepathNode) GetID() int {
	return m.ID
}

func (m *TelepathNode) UseID() bool {
	return m.UseIdentifier
}

func (m *TelepathNode) Emit() any {
	if m.UseID() {
		return TelepathValue{Ref: m.ID}
	}

	m.Seen = true
	if m.UseID() && m.ID != 0 {
		var result = m.EmitVerbose()
		result.ID = m.ID
		return result
	}

	return m.EmitCompact()
}

func (m *TelepathNode) EmitVerbose() TelepathValue {
	return TelepathValue{}
}

func (m *TelepathNode) EmitCompact() any {
	return nil
}

type TelepathValueNode struct {
	*TelepathNode
	Value interface{}
}

func NewTelepathValueNode(value interface{}) *TelepathValueNode {
	return &TelepathValueNode{
		Value:        value,
		TelepathNode: NewTelepathNode(),
	}
}

func (m *TelepathValueNode) UseID() bool {
	return false
}

func (m *TelepathValueNode) GetValue() interface{} {
	return m.Value
}

func (m *TelepathValueNode) EmitVerbose() TelepathValue {
	return TelepathValue{Val: m.GetValue()}
}

func (m *TelepathValueNode) EmitCompact() any {
	return m.GetValue()
}

type UUIDNode struct {
	*TelepathValueNode
}

func NewUUIDNode(value interface{}) *UUIDNode {
	var v, ok = value.(uuid.UUID)
	if !ok {
		panic("value is not of type google/uuid.UUID")
	}

	return &UUIDNode{
		TelepathValueNode: NewTelepathValueNode(v),
	}
}

func (m *UUIDNode) Emit() any {
	if m.UseID() {
		return TelepathValue{Ref: m.ID}
	}

	m.Seen = true
	if m.ID != 0 {
		var result = m.EmitVerbose()
		result.ID = m.ID
		return result
	}

	return m.EmitCompact()
}

func (m *UUIDNode) EmitVerbose() TelepathValue {
	return TelepathValue{Val: m.GetValue()}
}

func (m *UUIDNode) EmitCompact() any {
	return m.GetValue()
}

type StringNode struct {
	*TelepathValueNode
}

func (m *StringNode) UseID() bool {
	var rVal = reflect.ValueOf(m.GetValue())
	return rVal.Len() >= STRING_REF_MIN_LENGTH && m.ID != 0 && m.Seen
}

func NewStringNode(value interface{}) *StringNode {
	return &StringNode{
		TelepathValueNode: NewTelepathValueNode(value),
	}
}

func (m *StringNode) Emit() any {
	if m.UseID() && m.ID != 0 {
		return TelepathValue{Ref: m.ID}
	}

	m.Seen = true
	if m.ID != 0 {
		var result = m.EmitVerbose()
		result.ID = m.ID
		return result
	}

	return m.EmitCompact()
}

func (m *StringNode) EmitVerbose() TelepathValue {
	return TelepathValue{Val: m.GetValue()}
}

func (m *StringNode) EmitCompact() any {
	return m.GetValue()
}

type nullNode struct {
	*TelepathValueNode
}

func NullNode() *nullNode {
	return &nullNode{
		TelepathValueNode: NewTelepathValueNode(nil),
	}
}

type ErrorNode struct {
	*TelepathValueNode
}

func NewErrorNode(value error) *ErrorNode {
	return &ErrorNode{
		TelepathValueNode: NewTelepathValueNode(value),
	}
}

func (m *ErrorNode) Emit() any {
	if m.UseID() {
		return TelepathValue{Ref: m.ID}
	}

	m.Seen = true
	if m.ID != 0 {
		var result = m.EmitVerbose()
		result.ID = m.ID
		return result
	}

	return m.EmitCompact()
}

func (m *ErrorNode) EmitVerbose() TelepathValue {
	return TelepathValue{Val: m.GetValue().(error).Error()}
}

func (m *ErrorNode) EmitCompact() any {
	return m.GetValue().(error).Error()
}

type ObjectNode struct {
	*TelepathValueNode
	Constructor string
	Args        []Node
}

func NewObjectNode(constructor string, args []Node) *ObjectNode {
	return &ObjectNode{
		TelepathValueNode: NewTelepathValueNode(nil),
		Constructor:       constructor,
		Args:              args,
	}
}

func (m *ObjectNode) SetID(id int) {
	m.ID = id
	m.UseIdentifier = true
}

func (m *ObjectNode) UseID() bool {
	return m.ID != 0 && m.Seen
}

func (m *ObjectNode) Emit() any {
	if m.UseID() {
		return TelepathValue{Ref: m.ID}
	}

	m.Seen = true

	if m.ID != 0 {
		var result = m.EmitVerbose()
		result.ID = m.ID
		return result
	}

	return m.EmitCompact()
}

func (m *ObjectNode) EmitVerbose() TelepathValue {
	var result = TelepathValue{
		Type: m.Constructor,
		Args: make([]any, 0, len(m.Args)),
	}
	for _, arg := range m.Args {
		result.Args = append(result.Args, arg.Emit())
	}
	return result
}

func (m *ObjectNode) EmitCompact() any {
	return m.EmitVerbose()
}

type DictNode struct {
	*TelepathValueNode
}

func (m *DictNode) UseID() bool {
	return m.ID != 0 && m.UseIdentifier
}

func (m *DictNode) SetID(id int) {
	m.ID = id
	m.UseIdentifier = true
}
func NewDictNode(value map[string]Node) *DictNode {
	return &DictNode{
		TelepathValueNode: NewTelepathValueNode(value),
	}
}

func (m *DictNode) Emit() any {
	if m.UseID() {
		return TelepathValue{Ref: m.ID}
	}

	m.Seen = true
	if m.ID != 0 {
		var result = m.EmitVerbose()
		result.ID = m.ID
		return result
	}

	return m.EmitCompact()
}

func (m *DictNode) EmitVerbose() TelepathValue {
	var result = TelepathValue{Dict: make(map[string]interface{})}
	for key, value := range m.Value.(map[string]Node) {
		result.Dict[key] = value.Emit()
	}
	return result
}

func (m *DictNode) EmitCompact() any {
	var (
		hasReservedKey = false
		result         = make(map[string]interface{})
	)

	for key := range m.Value.(map[string]Node) {
		_, hasReservedKey = slices.BinarySearch(
			DICT_RESERVED_KEYS, key,
		)
		if hasReservedKey {
			return m.EmitVerbose()
		}
	}

	for key, value := range m.Value.(map[string]Node) {
		result[key] = value.Emit()
	}

	return result
}

type ListNode struct {
	*TelepathValueNode
}

func NewListNode(value []Node) *ListNode {
	return &ListNode{
		TelepathValueNode: NewTelepathValueNode(value),
	}
}

func (m *ListNode) GetValue() interface{} {
	return m.Value
}

func (m *ListNode) Emit() any {
	if m.UseID() {
		return TelepathValue{Ref: m.ID}
	}

	m.Seen = true

	var result = m.EmitVerbose()
	if m.ID != 0 {
		result.ID = m.ID
	}
	return result
}

func (m *ListNode) UseID() bool {
	return m.ID != 0 && m.Seen
}

func (m *ListNode) SetID(id int) {
	m.ID = id
	m.UseIdentifier = true
}

func (m *ListNode) EmitVerbose() TelepathValue {
	var result = TelepathValue{List: make([]interface{}, 0)}
	for _, value := range m.Value.([]Node) {
		result.List = append(result.List, value.Emit())
	}
	return result
}

func (m *ListNode) EmitCompact() any {
	var result = make([]interface{}, 0)
	for _, value := range m.Value.([]Node) {
		result = append(result, value.Emit())
	}
	return result
}
