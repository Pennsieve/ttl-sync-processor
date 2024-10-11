package fromcuration

import (
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/curation"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

// ToSample is a mappings.Mapping from curation.Sample to metadata.Sample
func ToSample(exportedSample curation.Sample) (metadata.Sample, error) {
	return metadata.Sample{ID: changesetmodels.ExternalInstanceID(exportedSample.ID)}, nil
}
