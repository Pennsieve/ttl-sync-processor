package models

// DatasetCurationExport represents the fields of interest from a single curation-export.json file for a dataset
// ID is the full Pennsieve dataset node id
type DatasetCurationExport struct {
	ID           string              `json:"id"`
	Contributors []Contributor       `json:"contributors"`
	DirStructure []DirStructureEntry `json:"dir_structure"`
	SpecimenDirs SpecimenDirs        `json:"specimen_dirs"`
}

func NewDatasetCurationExport(datasetID string) *DatasetCurationExport {
	return &DatasetCurationExport{ID: datasetID}
}

func (d *DatasetCurationExport) WithContributors(contributors ...Contributor) *DatasetCurationExport {
	d.Contributors = append(d.Contributors, contributors...)
	return d
}

func (d *DatasetCurationExport) WithDirStructureEntries(entry ...DirStructureEntry) *DatasetCurationExport {
	d.DirStructure = append(d.DirStructure, entry...)
	return d
}

func (d *DatasetCurationExport) WithSpecimenDirs(dirs SpecimenDirs) *DatasetCurationExport {
	d.SpecimenDirs = dirs
	return d
}
