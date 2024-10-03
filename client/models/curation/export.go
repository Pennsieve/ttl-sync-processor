package curation

// DatasetExport represents the fields of interest from a single curation-export.json file for a dataset
// ID is the full Pennsieve dataset node id
type DatasetExport struct {
	ID           string              `json:"id"`
	Contributors []Contributor       `json:"contributors"`
	DirStructure []DirStructureEntry `json:"dir_structure"`
	SpecimenDirs SpecimenDirs        `json:"specimen_dirs"`
	Subjects     []Subject           `json:"subjects"`
}

func NewDatasetExport(datasetID string) *DatasetExport {
	return &DatasetExport{ID: datasetID}
}

func (d *DatasetExport) WithContributors(contributors ...Contributor) *DatasetExport {
	d.Contributors = append(d.Contributors, contributors...)
	return d
}

func (d *DatasetExport) WithDirStructureEntries(entry ...DirStructureEntry) *DatasetExport {
	d.DirStructure = append(d.DirStructure, entry...)
	return d
}

func (d *DatasetExport) WithSpecimenDirs(dirs SpecimenDirs) *DatasetExport {
	d.SpecimenDirs = dirs
	return d
}

func (d *DatasetExport) WithSubjects(subjects ...Subject) *DatasetExport {
	d.Subjects = append(d.Subjects, subjects...)
	return d
}
