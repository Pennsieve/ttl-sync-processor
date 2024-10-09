package sync

import (
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/spec"
	"log/slog"
)

func ComputeIdentifiableModelChanges[OLD metadata.SavedIDer, NEW metadata.IDer](
	schemaData *metadataclient.Schema,
	old []OLD,
	new []NEW,
	instanceSpec spec.IdentifiableInstance[OLD, NEW]) (*changesetmodels.ModelChanges, error) {
	modelChanges, err := addIdentifiableModelChanges(old, new, instanceSpec)
	if err != nil {
		return nil, err
	}
	modelLogger := logger.With(slog.String("modelName", instanceSpec.Model.Name))
	if modelChanges == nil {
		modelLogger.Info("no changes")
		return nil, nil
	}
	if err := setModelIDOrCreate(modelChanges, schemaData, instanceSpec.Model); err != nil {
		return nil, err
	}
	deleteMessage := slog.Int("deleteCount", len(modelChanges.Records.Delete))
	if modelChanges.Records.DeleteAll {
		deleteMessage = slog.Int("deleteAllCount", len(old))
	}
	modelLogger.Info("change summary",
		deleteMessage,
		slog.Int("createCount", len(modelChanges.Records.Create)),
		slog.Int("updateCount", len(modelChanges.Records.Update)),
	)
	return modelChanges, nil
}
func addIdentifiableModelChanges[OLD metadata.SavedIDer, NEW metadata.IDer](old []OLD, new []NEW, instanceSpec spec.IdentifiableInstance[OLD, NEW]) (*changesetmodels.ModelChanges, error) {
	recordChanges := changesetmodels.RecordChanges{}

	oldByID := map[string]OLD{}
	oldToDelete := map[string]OLD{}
	for _, oldInstance := range old {
		oldByID[oldInstance.GetID()] = oldInstance
		oldToDelete[oldInstance.GetID()] = oldInstance
	}

	for _, newInstance := range new {
		newID := newInstance.GetID()
		if _, found := oldToDelete[newID]; found {
			delete(oldToDelete, newID)
		}
		if oldInstance, exists := oldByID[newID]; !exists {
			recordChanges.Create = append(recordChanges.Create, instanceSpec.Creator(newInstance))
		} else {
			recordUpdate, err := instanceSpec.Updater(oldInstance, newInstance)
			if err != nil {
				return nil, err
			}
			if recordUpdate != nil {
				recordChanges.Update = append(recordChanges.Update, *recordUpdate)
			}
		}
	}

	if len(oldToDelete) > 0 {
		if len(oldToDelete) == len(old) {
			// use batch delete if we're going to delete all the existing records anyway
			recordChanges.DeleteAll = true
		} else {
			for _, toDelete := range oldToDelete {
				recordChanges.Delete = append(recordChanges.Delete, toDelete.GetPennsieveID())
			}
		}
	}
	if recordChanges.DeleteAll == false && len(recordChanges.Create) == 0 && len(recordChanges.Delete) == 0 && len(recordChanges.Update) == 0 {
		return nil, nil
	}

	return &changesetmodels.ModelChanges{Records: recordChanges}, nil
}

func addIdentifiablePropertyLinkChanges[OLD metadata.SavedIDer, NEW metadata.IDer](old []OLD, new []NEW, instanceSpec spec.IdentifiableInstance[OLD, NEW]) (*changesetmodels.LinkedPropertyChanges, error) {
	instanceChanges := changesetmodels.InstanceChanges{}

	oldByID := map[string]OLD{}
	oldToDelete := map[string]OLD{}
	for _, oldInstance := range old {
		oldByID[oldInstance.GetID()] = oldInstance
		oldToDelete[oldInstance.GetID()] = oldInstance
	}

	for _, newInstance := range new {
		newID := newInstance.GetID()
		if _, found := oldToDelete[newID]; found {
			delete(oldToDelete, newID)
		}
		if _, exists := oldByID[newID]; !exists {
			instanceChanges.Creates = append(instanceChanges.Creates, changesetmodels.InstanceLinkedPropertyCreate{
				FromRecordID: "",
				InstanceLinkedPropertyCreatePayload: changesetmodels.InstanceLinkedPropertyCreatePayload{
					Name:                   "",
					DisplayName:            "",
					SchemaLinkedPropertyID: "",
					ToRecordID:             "",
				},
			})
		}
	}

	if len(oldToDelete) > 0 {
		for _, toDelete := range oldToDelete {
			instanceChanges.Deletes = append(instanceChanges.Deletes, changesetmodels.InstanceLinkedPropertyDelete{
				FromRecordID:             "",
				InstanceLinkedPropertyID: toDelete.GetPennsieveID(),
			})
		}
	}

	if len(instanceChanges.Creates) == 0 && len(instanceChanges.Deletes) == 0 {
		return nil, nil
	}
	return &changesetmodels.LinkedPropertyChanges{
		Instances: instanceChanges,
	}, nil
}
