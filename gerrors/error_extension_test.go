package gerrors

type ExtendedError struct {
	GError
	GRPCStatus int
}

var ErrTest2 Factory = &ExtendedError{
	GError:     GError{},
	GRPCStatus: 0,
}

var ErrTest Factory = &GError{
	Name:        "",
	Message:     "",
	Source:      "",
	UseFullSack: false,
	detailTag:   "",
	stack:       nil,
	srcFactory:  nil,
	srcError:    nil,
}
