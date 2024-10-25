package fromrecord

import (
	"fmt"
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-pre-metadata/client/models/instance"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/mappings"
)

func NewSampleStoreMapping(idMap *RecordIDStore) mappings.Mapping[instance.LinkedProperty, metadata.SavedSampleSubject] {
	return func(linkedProperty instance.LinkedProperty) (metadata.SavedSampleSubject, error) {
		savedSampleSubject := metadata.SavedSampleSubject{PennsieveID: changesetmodels.PennsieveInstanceID(linkedProperty.ID)}
		if err := checkLinkedPropertyName(linkedProperty, metadata.SampleSubjectLinkName); err != nil {
			return metadata.SavedSampleSubject{}, err
		}
		from := changesetmodels.PennsieveInstanceID(linkedProperty.From)
		to := changesetmodels.PennsieveInstanceID(linkedProperty.To)
		sampleSubject, err := GetSampleSubject(from, to, idMap)
		if err != nil {
			return metadata.SavedSampleSubject{}, err
		}
		savedSampleSubject.SampleSubject = sampleSubject
		savedSampleSubject.From = from
		savedSampleSubject.To = to
		return savedSampleSubject, nil
	}
}

func GetSampleSubject(sampleInstanceID, subjectInstanceID changesetmodels.PennsieveInstanceID, idMap *RecordIDStore) (metadata.SampleSubject, error) {
	sampleSubject := metadata.SampleSubject{}
	if sampleKey := idMap.GetExternal(sampleInstanceID); sampleKey == nil {
		return metadata.SampleSubject{}, fmt.Errorf("no sample with Pennsieve ID %s", sampleInstanceID)
	} else {
		sampleSubject.SampleID = sampleKey.ExternalRecordID
	}
	if subjectKey := idMap.GetExternal(subjectInstanceID); subjectKey == nil {
		return metadata.SampleSubject{}, fmt.Errorf("no subject with Pennsieve ID %s", subjectInstanceID)

	} else {
		sampleSubject.SubjectID = subjectKey.ExternalRecordID
	}
	return sampleSubject, nil
}
