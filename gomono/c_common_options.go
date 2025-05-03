package gomono

// ParentCommitOpt is a reusable type that can be embedded in any command to get the flag value.
type ParentCommitOpt struct {
	ParentCommit string `long:"parent" description:"The identifier of a parent commit (main, head, commit hash, etc) defaults to ALL modules if not specified."`
}
