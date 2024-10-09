package sync

import (
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/spec"
)

func ComputeSubjectChanges(schemaData *metadataclient.Schema, old []metadata.SavedSubject, new []metadata.Subject) (*changesetmodels.ModelChanges, error) {
	return ComputeIdentifiableModelChanges(schemaData, old, new, spec.SubjectInstance)
}
