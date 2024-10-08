package sync

import (
	"github.com/pennsieve/processor-pre-metadata/client/models/schema"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/spec"
)

func ComputeSampleChanges(schemaData map[string]schema.Element, old []metadata.SavedSample, new []metadata.Sample) (*changesetmodels.ModelChanges, error) {
	return ComputeIdentifiableModelChanges(schemaData, old, new, spec.SampleInstance)
}
