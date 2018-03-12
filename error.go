package trompe

import "fmt"

const (
	GenericError = iota
	InvalidArityError
	KeyError
)

type RuntimeError struct {
	Context *Context
	Type    int
	Reason  string
}

func ErrorName(ty int) string {
	switch ty {
	case GenericError:
		return "GenericError"
	case InvalidArityError:
		return "InvalidArityError"
	default:
		panic("unknown error")
	}
}

func CreateRuntimeError(ctx *Context, ty int, reason string) *RuntimeError {
	return &RuntimeError{ctx, ty, reason}
}

func (err *RuntimeError) Error() string {
	return fmt.Sprintf("%s: %s", ErrorName(err.Type), err.Reason)
}

func CreateInvalidArityError(ctx *Context, nargs int) *RuntimeError {
	return CreateRuntimeError(ctx, InvalidArityError, "")
}

func CreateKeyError(ctx *Context, name string) *RuntimeError {
	return CreateRuntimeError(ctx, KeyError,
		fmt.Sprintf("key %s not found", name))
}

func ValidateArity(ctx *Context, expected int, actual int) *RuntimeError {
	if expected != actual {
		return CreateRuntimeError(ctx,
			InvalidArityError,
			fmt.Sprintf("invalid arity (takes %d, but %d given)", expected, actual))
	} else {
		return nil
	}
}
