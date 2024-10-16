package fromrecord

import (
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-pre-metadata/client/models/instance"
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
			sample.ID = safeExternalID(v.Value)
		case metadata.PrimaryKeyKey:
			sample.PrimaryKey = safeString(v.Value)
		}
	}
	return sample, nil
}
