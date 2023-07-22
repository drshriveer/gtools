package error

type GError interface {
}

type GErrorImpl struct {
	Source   string
	GRPCCode codes.Code
}
