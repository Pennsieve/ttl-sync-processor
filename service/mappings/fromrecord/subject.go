package fromrecord

import (
	"github.com/pennsieve/processor-pre-metadata/client/models/instance"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

func ToSubject(record instance.Record) (metadata.SavedSubject, error) {
	subject := metadata.SavedSubject{PennsieveID: changesetmodels.PennsieveInstanceID(record.ID)}
	if err := checkRecordType(record, metadata.SubjectModelName); err != nil {
		return subject, err
	}
	for _, v := range record.Values {
		switch v.Name {
		case metadata.SubjectIDKey:
			subject.ID = safeExternalID(v.Value)
		case metadata.SpeciesKey:
			subject.Species = safeString(v.Value)
		case metadata.SpeciesSynonymsKey:
			subject.SpeciesSynonyms = safeStringSlice(v.Value)
		}
	}
	return subject, nil

}
