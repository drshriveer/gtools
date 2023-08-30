package errors

// Source represents the source of an error in any form it may take.
type Source interface {
	// Metric returns a single metric-safe string indicating source.
	Metric() string

	// TODO: not 100% sure i want to do this...
	isSource()
}

// External represents a source outside of known packages.
type External struct {
	Service string
	RPC     string
	Detail  string
}

func (d *External) Metric() string {
	return convertToMetricNode(d.Service, d.RPC, d.Detail)
}

func (External) isSource() {}
