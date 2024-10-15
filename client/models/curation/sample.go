package curation

type Sample struct {
	ID        string `json:"sample_id"`
	SubjectID string `json:"subject_id"`
	// PrimaryKey is the join between a Sample and the SpecimenID in a SpecimenDirs
	// So we join on Sample.PrimaryKey == SpecimenDirs.Records[x].SpecimenID when SpecimenDirs.Records[x].Type == SampleRecordType
	PrimaryKey string `json:"primary_key"`
}
