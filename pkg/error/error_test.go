package error_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

type AThing struct {
}

func (a *AThing) WithAMethod(t *testing.T) string {
	return makeStack(t)
}

func TestStack(t *testing.T) {
	t.Run("subtask", func(t *testing.T) {
		L1(t)
		thing := AThing{}
		println(thing.WithAMethod(t))
	},
	)
}

func L1(t *testing.T) {
	L2(t)
}
func L2(t *testing.T) {
	L3(t)
}
func L3(t *testing.T) {
	println(makeStack(t))
}

type StackElem struct {
	Name string
	File string
	Line int
}

// TODO: can I find the package name of where a defined error originated?
// through reflect? if the item is a pointer to the original?
// omg that could be used as the "NAME"
//

func (e StackElem) String() string {
	splitName := strings.Split(e.Name, "/")
	last := splitName[len(splitName)-1]
	vals := strings.Split(last, ".")

	var packageName, funcName, typeName string
	if len(vals) == 2 {
		packageName, funcName = vals[0], vals[1]
	} else if len(vals) == 3 {
		packageName, typeName, funcName = vals[0], vals[1], vals[2]
		typeName = strings.TrimPrefix(typeName, "(")
		typeName = strings.TrimSuffix(typeName, ")")
	}

	return fmt.Sprintf("[%s, %s, %s] ",
		packageName,
		typeName,
		funcName,
	) + e.Name + " " + e.File + ":" + strconv.Itoa(e.Line)
}

type Stack []StackElem

func (s Stack) String() string {
	sb := strings.Builder{}
	for _, e := range s {
		sb.WriteString("\n\t" + e.String())
	}
	return sb.String()
}

func makeStack(t *testing.T) string {
	depth := 32
	pcs := make([]uintptr, depth)
	n := runtime.Callers(2, pcs[:])
	pcs = pcs[0:n] // drop unwritten elements.
	assert.Len(t, pcs, n)
	stack := make(Stack, n)
	for i := range stack {
		// program counter for line
		pc := pcs[i] - 1
		fu := runtime.FuncForPC(pc)
		if fu == nil {
			stack[i] = StackElem{Name: "unknown", File: "unknown"}
		} else {
			fName, fLine := fu.FileLine(pc)
			fu.Name()
			stack[i] = StackElem{Name: fu.Name(), File: fName, Line: fLine}
		}
	}

	return stack.String()
}
