package fromrecord

import (
	"fmt"
	"github.com/pennsieve/processor-pre-metadata/client/models/instance"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/mappings"
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

func NewSampleStoreMapping(samples []metadata.SavedSample, subjects []metadata.SavedSubject) mappings.Mapping[instance.LinkedProperty, metadata.SavedSampleSubjectInstance] {
	store := NewSampleSubjectStore(samples, subjects)
	return func(linkedProperty instance.LinkedProperty) (metadata.SavedSampleSubjectInstance, error) {
		savedSampleSubject := metadata.SavedSampleSubjectInstance{PennsieveID: changesetmodels.PennsieveInstanceID(linkedProperty.Id)}
		if err := checkLinkedPropertyName(linkedProperty, metadata.SampleSubjectLinkName); err != nil {
			return metadata.SavedSampleSubjectInstance{}, err
		}
		from := changesetmodels.PennsieveInstanceID(linkedProperty.From)
		to := changesetmodels.PennsieveInstanceID(linkedProperty.To)
		sampleSubject, err := store.GetSampleSubject(from, to)
		if err != nil {
			return metadata.SavedSampleSubjectInstance{}, err
		}
		savedSampleSubject.SampleSubject = sampleSubject
		savedSampleSubject.From = from
		savedSampleSubject.To = to
		return savedSampleSubject, nil
	}
}

type SampleSubjectStore struct {
	samplesByPennsieveID  map[changesetmodels.PennsieveInstanceID]metadata.SavedSample
	subjectsByPennsieveID map[changesetmodels.PennsieveInstanceID]metadata.SavedSubject
}

func NewSampleSubjectStore(samples []metadata.SavedSample, subjects []metadata.SavedSubject) *SampleSubjectStore {
	store := &SampleSubjectStore{}
	store.samplesByPennsieveID = make(map[changesetmodels.PennsieveInstanceID]metadata.SavedSample, len(samples))
	for _, sample := range samples {
		store.samplesByPennsieveID[sample.GetPennsieveID()] = sample
	}
	store.subjectsByPennsieveID = make(map[changesetmodels.PennsieveInstanceID]metadata.SavedSubject, len(subjects))
	for _, subject := range subjects {
		store.subjectsByPennsieveID[subject.GetPennsieveID()] = subject
	}
	return store
}

func (store *SampleSubjectStore) GetSampleSubject(sampleInstanceID, subjectInstanceID changesetmodels.PennsieveInstanceID) (metadata.SampleSubject, error) {
	sampleSubject := metadata.SampleSubject{}
	if sample, sampleFound := store.samplesByPennsieveID[sampleInstanceID]; !sampleFound {
		return metadata.SampleSubject{}, fmt.Errorf("no sample with Pennsieve ID %s", sampleInstanceID)
	} else {
		sampleSubject.SampleID = sample.GetID()
	}
	if subject, subjectFound := store.subjectsByPennsieveID[subjectInstanceID]; !subjectFound {
		return metadata.SampleSubject{}, fmt.Errorf("no subject with Pennsieve ID %s", subjectInstanceID)

	} else {
		sampleSubject.SubjectID = subject.GetID()
	}
	return sampleSubject, nil
}
