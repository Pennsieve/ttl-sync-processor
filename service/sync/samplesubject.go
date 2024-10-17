package sync

import (
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/spec"
)

func ComputeSampleSubjectChanges(schemaData *metadataclient.Schema, old []metadata.SavedSampleSubject, new []metadata.SampleSubject) (*changesetmodels.LinkedPropertyChanges, error) {
	return ComputeIdentifiableLinkedPropertyChanges(schemaData, old, new, spec.SampleSubject)
}
