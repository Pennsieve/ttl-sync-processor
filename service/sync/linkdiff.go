package sync

import (
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/spec"
	"log/slog"
)

func ComputeIdentifiableLinkedPropertyChanges[OLD metadata.SavedExternalLink, NEW metadata.ExternalLink](
	schemaData *metadataclient.Schema,
	old []OLD,
	new []NEW,
	linkSpec spec.Link) (*changesetmodels.LinkedPropertyChanges, error) {
	linkChanges, err := addIdentifiableLinkedPropertyChanges(old, new)
	if err != nil {
		return nil, err
	}
	linkLogger := logger.With(slog.String("linkName", linkSpec.Name))
	if linkChanges == nil {
		linkLogger.Info("no changes")
		return nil, nil
	}
	if err := setLinkIDOrCreate(linkChanges, schemaData, linkSpec); err != nil {
		return nil, err
	}
	linkLogger.Info("change summary",
		slog.Int("createCount", len(linkChanges.Instances.Create)),
		slog.Int("deleteCount", len(linkChanges.Instances.Delete)),
	)
	return linkChanges, nil
}
func addIdentifiableLinkedPropertyChanges[OLD metadata.SavedExternalLink, NEW metadata.ExternalLink](old []OLD, new []NEW) (*changesetmodels.LinkedPropertyChanges, error) {
	instanceChanges := changesetmodels.InstanceChanges{}

	oldByID := map[changesetmodels.ExternalInstanceID]OLD{}
	oldToDelete := map[changesetmodels.ExternalInstanceID]OLD{}
	for _, oldInstance := range old {
		oldByID[oldInstance.ExternalID()] = oldInstance
		oldToDelete[oldInstance.ExternalID()] = oldInstance
	}

	for _, newInstance := range new {
		newID := newInstance.ExternalID()
		if _, found := oldToDelete[newID]; found {
			delete(oldToDelete, newID)
		}
		if _, exists := oldByID[newID]; !exists {
			instanceChanges.Create = append(instanceChanges.Create, changesetmodels.InstanceLinkedPropertyCreate{
				FromExternalID: newInstance.FromExternalID(),
				ToExternalID:   newInstance.ToExternalID(),
			})
		}
	}

	if len(oldToDelete) > 0 {
		for _, toDelete := range oldToDelete {
			instanceChanges.Delete = append(instanceChanges.Delete, changesetmodels.InstanceLinkedPropertyDelete{
				FromRecordID:             toDelete.FromPennsieveID(),
				InstanceLinkedPropertyID: toDelete.GetPennsieveID(),
			})
		}
	}

	if len(instanceChanges.Create) == 0 && len(instanceChanges.Delete) == 0 {
		return nil, nil
	}
	return &changesetmodels.LinkedPropertyChanges{
		Instances: instanceChanges,
	}, nil
}

func setLinkIDOrCreate(linkChanges *changesetmodels.LinkedPropertyChanges, schemaData *metadataclient.Schema, linkSpec spec.Link) error {
	if linkSchema, linkSchemaExists := schemaData.LinkedPropertyByName(linkSpec.Name); linkSchemaExists {
		logger.Info("linkSchema exists", slog.String("linkSchemaName", linkSchema.Name), slog.String("linkSchemaID", linkSchema.ID))
		linkChanges.ID = linkSchema.ID
	} else {
		schemaCreate := linkSpec.SchemaCreate()
		logger.Info("linkSchema must be created", slog.String("linkName", linkSpec.Name))
		linkChanges.Create = &schemaCreate
	}

	return nil
}
