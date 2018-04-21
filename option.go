package trompe

import (
	"fmt"
)

type Option struct {
	Value Value // nullable
}

func NewOption(v Value) *Option {
	return &Option{v}
}

func (o *Option) Type() int {
	return ValueTypeOption
}

func (o *Option) Desc() string {
	return fmt.Sprintf("<option %s>", o.Value.Desc())
}
