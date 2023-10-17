package gencommon

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Comments is a line-by-line representation of comments.
type Comments []string

// CommentsFromMethod extracts comments of a method with type and method name.
func CommentsFromMethod(pkg *packages.Package, typeName string, methodName string) Comments {
	for _, stax := range pkg.Syntax {
		obj := stax.Scope.Lookup(typeName)
		if obj != nil {
			// if this is an interface type extracting comments is easier.
			if v, ok := obj.Type.(*ast.InterfaceType); ok {
				for _, mInfo := range v.Methods.List {
					if getName(mInfo.Names...) == methodName {
						return FromCommentGroup(mInfo.Comment)
					}
				}
			}
		}

		// if this is a struct type we first have to find functions with receivers of the correct type.
		for _, decl := range stax.Decls {
			if fDecl, ok := decl.(*ast.FuncDecl); ok && methodIdent(fDecl, typeName) {
				if getName(fDecl.Name) == methodName {
					return FromCommentGroup(fDecl.Doc)
				}
			}
		}
	}

	return nil
}

// FromCommentGroup converts a CommentGroup into Comments!
func FromCommentGroup(group *ast.CommentGroup) Comments {
	if group == nil {
		return nil
	}
	return mapper(group.List, func(in *ast.Comment) string { return in.Text })
}

// String returns a single string of the comment block.
func (c Comments) String() string {
	return strings.Join(c, "\n")
}
