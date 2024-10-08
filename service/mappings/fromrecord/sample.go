package fromrecord

import (
	"github.com/pennsieve/processor-pre-metadata/client/models/instance"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

func ToSample(record instance.Record) (metadata.SavedSample, error) {
	sample := metadata.SavedSample{PennsieveID: changesetmodels.PennsieveInstanceID(record.ID)}
	if err := checkRecordType(record, metadata.SampleModelName); err != nil {
		return metadata.SavedSample{}, nil
	}
	for _, v := range record.Values {
		switch v.Name {
		case metadata.SampleIDKey:
			sample.ID = safeString(v.Value)
		}
	}
	return sample, nil
}

func ToSampleSubjectLink(linkedProperty instance.LinkedProperty) (metadata.SavedSampleSubjectLink, error) {
	sampleSubject := metadata.SavedSampleSubjectLink{PennsieveID: changesetmodels.PennsieveInstanceID(linkedProperty.Id)}
	if err := checkLinkedPropertyName(linkedProperty, metadata.SampleSubjectLinkName); err != nil {
		return metadata.SavedSampleSubjectLink{}, err
	}
	sampleSubject.SampleID = linkedProperty.From
	sampleSubject.SubjectID = linkedProperty.To
	return sampleSubject, nil
}
