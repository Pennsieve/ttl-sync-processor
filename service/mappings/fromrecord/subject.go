package fromrecord

import (
	"github.com/pennsieve/processor-pre-metadata/client/models/instance"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

func ToSubject(record instance.Record) (metadata.Subject, error) {
	subject := metadata.Subject{}
	if err := checkRecordType(record, metadata.SubjectModelName); err != nil {
		return subject, err
	}
	for _, v := range record.Values {
		switch v.Name {
		case metadata.SubjectIDKey:
			subject.ID = safeString(v.Value)
		case metadata.SpeciesKey:
			subject.Species = safeString(v.Value)
		case metadata.SpeciesSynonymsKey:
			subject.SpeciesSynonyms = safeStringSlice(v.Value)
		}
	}
	return subject, nil

}
