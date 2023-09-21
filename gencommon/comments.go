package gencommon

import (
	"go/ast"
)

// CommentGroupRaw returns the comment group as a single string
// but without manipulation.
func CommentGroupRaw(cg *ast.CommentGroup) string {
	result := ""
	for i, comment := range cg.List {
		result += comment.Text
		if i < len(cg.List)-1 {
			result += "\n"
		}
	}
	return result
}
