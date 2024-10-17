package curation

type SpecimenDirs struct {
	Records []Record `json:"records,omitempty"`
}

func NewSpecimenDirs() *SpecimenDirs {
	return &SpecimenDirs{}
}

func (d *SpecimenDirs) WithRecord(specimenID string, recordType RecordType, dir ...string) *SpecimenDirs {
	d.Records = append(d.Records, NewRecord(specimenID, recordType, dir...))
	return d
}

type RecordType string

const (
	SubjectRecordType RecordType = "SubjectDirs"
	SampleRecordType  RecordType = "SampleDirs"
)

type Record struct {
	Dirs       []string   `json:"dirs,omitempty"`
	SpecimenID string     `json:"specimen_id"`
	Type       RecordType `json:"type"`
}

func NewRecord(specimenID string, recordType RecordType, dir ...string) Record {
	return Record{
		Dirs:       dir,
		SpecimenID: specimenID,
		Type:       recordType,
	}
}
