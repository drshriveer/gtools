package gerror

// TODO: reconsider this delimiter when support for it is clearer.
const nonNodeMetricDelimiter = ":"

// convertToMetricNode takes a list of string elements and transforms
// them into a delimited metric-safe string skipping any empty entries.
func convertToMetricNode(elements ...string) string {
	result := ""
	for _, elem := range elements {
		if len(result) == 0 {
			result = elem
		} else if len(elem) > 0 {
			result += nonNodeMetricDelimiter + elem
		}
	}
	return result
}
