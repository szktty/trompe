package trompe

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type ObjectFile struct {
	Name     string        `json:"name"`
	Attrs    []*ObjectAttr `json:"attrs"`
	Codes    []*ObjectCode `json:"codes"`
	CodeVals map[int]*CompiledCode
	module   *Module
}

type ObjectAttr struct {
	Name  string       `json:"name"`
	Value *ObjectValue `json:"value"`
}

type ObjectValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type ObjectCode struct {
	Id   int            `json:"id"`
	Syms []string       `json:"symbols"`
	Lits []*ObjectValue `json:"literals"`
	Ops  []int          `json:"opcodes"`
}

var ObjectValueTypeUnit = "unit"
var ObjectValueTypeBool = "bool"
var ObjectValueTypeInt = "int"
var ObjectValueTypeFloat = "float"
var ObjectValueTypeString = "string"
var ObjectValueTypeCode = "code"

func NewObjectFile(name string) *ObjectFile {
	return &ObjectFile{
		Name:     name,
		Attrs:    []*ObjectAttr{},
		Codes:    []*ObjectCode{},
		CodeVals: make(map[int]*CompiledCode, 0),
	}
}

func (file *ObjectFile) AddAttr(attr *ObjectAttr) {
	file.Attrs = append(file.Attrs, attr)
}

func (file *ObjectFile) AddCode(code *ObjectCode) {
	file.Codes = append(file.Codes, code)
}

func (file *ObjectFile) AddCompiledCode(code *CompiledCode) {
	objCode := NewObjectCode(code.Id, code.Ops)
	for _, lit := range code.Lits {
		// TODO: other types
		switch lit.Type() {
		case ValClosType:
			if litCode, ok := lit.Closure().(*CompiledCode); ok {
				file.AddCompiledCode(litCode)
			} else {
				panic("notimpl")
			}
		default:
			break
		}
	}
	file.AddCode(objCode)
}

func NewObjectAttr(name string, value *ObjectValue) *ObjectAttr {
	return &ObjectAttr{Name: name, Value: value}
}

func NewObjectCode(id int, ops []Opcode) *ObjectCode {
	return &ObjectCode{Id: id, Syms: []string{}, Lits: []*ObjectValue{}, Ops: ops}
}

func (code *ObjectCode) AddSym(name string) {
	code.Syms = append(code.Syms, name)
}

func (code *ObjectCode) AddLit(value *ObjectValue) {
	code.Lits = append(code.Lits, value)
}

func (code *ObjectCode) AddOp(op int) {
	code.Ops = append(code.Ops, op)
}

func NewObjectValue(ty string, value string) *ObjectValue {
	return &ObjectValue{Type: ty, Value: value}
}

func NewObjectValueBool(value bool) *ObjectValue {
	var s string
	if value {
		s = "true"
	} else {
		s = "false"
	}
	return NewObjectValue(ObjectValueTypeBool, s)
}

func NewObjectValueCode(i int) *ObjectValue {
	return NewObjectValue(ObjectValueTypeCode, fmt.Sprintf("%s", i))
}

// Marshal/Unmarshal

func (file *ObjectFile) Marshal() ([]byte, error) {
	return json.Marshal(file)
}

func UnmarshalObjectFile(data []byte) (*ObjectFile, error) {
	var objFile ObjectFile
	if err := json.Unmarshal(data, &objFile); err != nil {
		return nil, err
	} else {
		return &objFile, nil
	}
}

func (file *ObjectFile) Decode() *Module {
	if file.module != nil {
		return file.module
	}

	m := NewModule(nil, file.Name)
	for _, objCode := range file.Codes {
		objCode.Decode(file)
	}
	for _, objAttr := range file.Attrs {
		m.AddAttr(objAttr.Name, objAttr.Value.Decode(file))
	}
	file.module = m
	return m
}

func (objCode *ObjectCode) Decode(file *ObjectFile) {
	code := NewCompiledCode()
	code.Id = objCode.Id
	code.Syms = objCode.Syms
	code.Ops = objCode.Ops
	file.CodeVals[code.Id] = code

	for _, objVal := range objCode.Lits {
		value := objVal.Decode(file)
		code.AddLit(value)
	}
}

func (value *ObjectValue) Decode(file *ObjectFile) Value {
	switch value.Type {
	case ObjectValueTypeUnit:
		return LangUnit
	case ObjectValueTypeString:
		return NewValStr(value.Value)
	case ObjectValueTypeCode:
		// TODO: error
		id, _ := strconv.Atoi(value.Value)
		return file.CodeVals[id]
	default:
		panic("notimpl")
	}
}
