package fromcuration

// Mapping is a function that converts an object from SPARC's curation export file (type From)
// to one of the Metadata model types defined for this sync (type To)
type Mapping[From, To any] func(curationObject From) (To, error)

// MapSlice uses the given Mapping to map a slice of Froms to a slice of Tos
func MapSlice[From, To any](exported []From, mapping Mapping[From, To]) ([]To, error) {
	var pennsieveMetadata []To
	for _, e := range exported {
		m, err := mapping(e)
		if err != nil {
			return nil, err
		}
		pennsieveMetadata = append(pennsieveMetadata, m)
	}
	return pennsieveMetadata, nil
}
