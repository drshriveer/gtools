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
			v, ok := obj.Decl.(*ast.InterfaceType)
			if !ok {
				if ts, tsok := obj.Decl.(*ast.TypeSpec); tsok {
					v, ok = ts.Type.(*ast.InterfaceType)
				}
			}
			// if this is an interface type extracting comments is easier.
			if ok {
				for _, mInfo := range v.Methods.List {
					if getName(mInfo.Names...) == methodName {
						return FromCommentGroup(mInfo.Doc)
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

// CommentsFromObj attempts to find comments associated with the typename provided in
// the package provided.
func CommentsFromObj(pkg *packages.Package, typeName string) Comments {
	for _, stax := range pkg.Syntax {
		obj := stax.Scope.Lookup(typeName)
		if obj == nil {
			continue
		}

		if decl, ok := obj.Decl.(*ast.TypeSpec); ok {
			// structs/interfaces won't have comments built-in unless
			// they're defined a particular way... e.g.
			// type {
			//    // comment on MyType
			//    MyType struct {
			//    ...
			//    }
			// }
			// So we'll check that first, but then build a comment map and try to get the
			// correct comment that way.
			if decl.Doc != nil && len(decl.Doc.List) > 0 {
				return FromCommentGroup(decl.Doc)
			} else if decl.Comment != nil && len(decl.Comment.List) > 0 {
				return FromCommentGroup(decl.Comment)
			}
			cmap := ast.NewCommentMap(pkg.Fset, stax, stax.Comments)
			for tp, commentCommentGroups := range cmap {
				if v, ok := tp.(*ast.GenDecl); ok && len(v.Specs) == 1 && v.Specs[0] == decl && len(commentCommentGroups) > 0 {
					// take the LAST comment.... theoretically the closest comment to the actual struct we're looking at.
					return FromCommentGroup(commentCommentGroups[len(commentCommentGroups)-1])
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
