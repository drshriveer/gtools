package gerrors

import (
	"runtime"
	"strconv"
	"strings"
)

// A Stack represents each line of a Stack trace.
type Stack []StackElem

// String returns a formatted sting with all
func (s Stack) String() string {
	sb := strings.Builder{}
	for _, e := range s {
		sb.WriteString("\n" + e.Name)
		sb.WriteString("\n\t" + e.File + ":" + strconv.Itoa(e.Line))
	}
	return sb.String()
}

// StackElem represents a single line in a Stack trace.
type StackElem struct {
	Name string
	File string
	Line int
}

// Source returns info about where this error was porpegated including packageName, typeName, and functionName.
func (e StackElem) Source() (packageName string, typeName string, funcName string) {
	splitName := strings.Split(e.Name, "/")
	last := splitName[len(splitName)-1]
	// Next step: handle generics which show up as funcName[....]
	// I'd love to do:
	//    last = strings.Replace(last, "[...]", "[T]", 1)
	// but this probably isn't metric safe.
	// I also assume that [...] handles N types so [T] wouldn't quite work.
	last = strings.TrimSuffix(last, "[...]")
	vals := strings.Split(last, ".")

	if len(vals) == 2 {
		packageName, funcName = vals[0], vals[1]
	} else if len(vals) == 3 {
		packageName, typeName, funcName = vals[0], vals[1], vals[2]
		typeName = strings.TrimPrefix(typeName, "(")
		typeName = strings.TrimSuffix(typeName, ")")
	}

	return packageName, typeName, funcName
}

// Metric returns a metric-safe(?) string of the source info.
func (e StackElem) Metric() string {
	pkg, tName, fName := e.Source()
	tName = strings.TrimPrefix(tName, "*") // remove pointer indicator
	return convertToMetricNode(pkg, tName, fName)
}

func (StackElem) isSource() {}

func makeStack(depth, skip int) *Stack {
	pcs := make([]uintptr, depth)
	n := runtime.Callers(skip, pcs[:])
	pcs = pcs[0:n] // drop unwritten elements.
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

	return &stack
}
