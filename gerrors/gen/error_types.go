package gen

import (
	"go/types"
)

// ErrorDesc describes an error.
//
//go:gen gensort --types=ErrorDesc
type ErrorDesc struct {
	TypeName string
	Type     types.Type
}
