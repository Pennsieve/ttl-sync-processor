package mappings

// Mapping is a function that given an object of type From
// returns a new object of type To
type Mapping[From, To any] func(curationObject From) (To, error)

// MapSlice uses the given Mapping to map a []From to a []To
func MapSlice[From, To any](source []From, mapping Mapping[From, To]) ([]To, error) {
	var tos []To
	for _, e := range source {
		m, err := mapping(e)
		if err != nil {
			return nil, err
		}
		tos = append(tos, m)
	}
	return tos, nil
}
