package sync

import (
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/spec"
)

func ComputeSampleChanges(schemaData *metadataclient.Schema, old []metadata.SavedSample, new []metadata.Sample) (any, error) {
	return ComputeIdentifiableModelChanges(schemaData, old, new, spec.SampleInstance)
}
