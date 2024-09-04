package models

// DatasetCurationExport represents the fields of interest from a single curation-export.json file for a dataset
// ID is the full Pennsieve dataset node id
type DatasetCurationExport struct {
	ID           string        `json:"id"`
	Contributors []Contributor `json:"contributors"`
}
