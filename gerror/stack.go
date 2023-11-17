package gerror

import (
	"runtime"
	"strconv"
	"strings"
)

// StackType identifies the depth of a stack desired.
// Generating stacks requires significant computation, smaller stacks use less.
type StackType int

const (
	// NoStack means do not generate a stack.
	NoStack StackType = 0
	// SourceStack retrieves the minimum sack possible to populate "source".
	// Note; this is a noop for errors with a defined source.
	// CAUTION: this is dubious... We're adding multiple stack elements to try and find
	// a more valid stack.... who knows if it will work.
	SourceStack StackType = 4
	// ShortStack gets a max stack of 16 elements.
	ShortStack StackType = 16
	// DefaultStack gets a max stack of 32 elements.
	DefaultStack StackType = 32
)

// StackSkip indicates how many stack layers to skip to get the correct start point.
type StackSkip int

const (
	// defaultSkip is 4 because that is how many layers the stack processor itself consumes.
	defaultSkip = 4
)

// A Stack represents each line of a Stack trace.
type Stack []StackElem

// String returns a formatted sting with all.
func (s Stack) String() string {
	sb := strings.Builder{}
	for _, e := range s {
		sb.WriteString("\n" + e.Name)
		sb.WriteString("\n\t" + e.File + ":" + strconv.Itoa(e.LineNumber))
	}
	return sb.String()
}

// NearestExternal finds the firs caller outside this package.
// The effectiveness of this method is limited to the depth of the Stack fetched.
func (s Stack) NearestExternal() StackElem {
	// attempt to find the first element not in this package:
	if pkgName, ok := getCurrentPackage(); ok {
		for _, elem := range s {
			if !strings.HasPrefix(elem.Name, pkgName) {
				return elem
			}
		}
	}
	return s[0]
}

func getCurrentPackage() (string, bool) {
	pc, _, _, _ := runtime.Caller(1)
	splitName := strings.Split(pcToStackElem(pc).Name, ".")
	if len(splitName) == 0 { // should literally be impossible?
		return "", false
	}
	return strings.Join(splitName[:len(splitName)-1], "."), true
}

// StackElem represents a single line in a Stack trace.
type StackElem struct {
	// Name is the fully qualified package function path (?).
	// e.g. github.com/drshriveer/gtools/gerror.TestGError_WithStack
	Name string

	// File is the full path of the file.
	File string

	// LineNumber of the Stack element.
	LineNumber int
}

// SourceInfo returns info about where this error was propagated including packageName,
// and other identifying information.
// The ambiguity of the other identifying information is unfortunate; but that's what we're
// starting with. I'd love info here...
func (e StackElem) SourceInfo() (packageName string, parts []string) {
	splitName := strings.Split(e.Name, "/")
	last := splitName[len(splitName)-1]
	// Next step: handle generics which show up as funcName[....]
	// I'd love to do:
	//    last = strings.Replace(last, "[...]", "[T]", 1)
	// but this probably isn't metric safe.
	// I also assume that [...] handles N types so maybe [T] wouldn't quite work.
	// FIXME: very possible the [...] only comes out in specific situaions.
	last = strings.TrimSuffix(last, "[...]")

	// "last" has a number of variations, I don't yet understand the rules.
	// (See tests for examples / more info.)
	// But generally, the first 3 nodes seem to be the most useful.
	// the first of which is _always_ a package.
	vals := strings.Split(last, ".")

	if len(vals) >= 1 {
		packageName = vals[0]
		vals = vals[1:]
	}

	// cut it down-- anything after funcX is not very useful imo
	for i, val := range vals {
		// conditions for truncaing the set.
		if strings.HasPrefix(val, "func") || val == "" {
			vals = vals[:i]
			break
		}
	}

	return packageName, vals
}

// Metric returns a metric-safe(?) string of the source info.
func (e StackElem) Metric() string {
	// delim is the metric node delimiter.
	// TODO: reconsider this delimiter when support for it is clearer.
	const delim = ":"

	pkg, theRest := e.SourceInfo()

	// Experimental .. drop repeats.
outer:
	for i := 1; i < len(theRest); i++ {
		curr := theRest[i]
		for j := 0; j < i; j++ {
			if curr == theRest[j] {
				// only use everything before the repeat.
				theRest = theRest[:i]
				break outer
			}
		}
	}

	// convertToMetricNode takes a list of string elements and transforms
	// them into a delimited metric-safe string skipping any empty entries.
	return pkg + delim + strings.Join(theRest, delim)

}

func makeStack(depth StackType, skip StackSkip) Stack {
	pcs := make([]uintptr, depth)
	n := runtime.Callers(int(skip), pcs)
	pcs = pcs[0:n] // drop unwritten elements.
	stack := make(Stack, n)
	for i := range stack {
		stack[i] = pcToStackElem(pcs[i])
	}

	return stack
}

func pcToStackElem(pc uintptr) StackElem {
	pc--
	fu := runtime.FuncForPC(pc)
	if fu == nil {
		return StackElem{Name: "unknown", File: "unknown"}
	}
	fName, fLine := fu.FileLine(pc)
	fu.Entry()
	return StackElem{Name: fu.Name(), File: fName, LineNumber: fLine}
}
