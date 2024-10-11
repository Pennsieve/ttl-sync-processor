package sync

import (
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/spec"
	"log/slog"
)

func ComputeIdentifiableModelChanges[OLD metadata.SavedExternalIDer, NEW metadata.ExternalIDer](
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
func addIdentifiableModelChanges[OLD metadata.SavedExternalIDer, NEW metadata.ExternalIDer](old []OLD, new []NEW, instanceSpec spec.IdentifiableInstance[OLD, NEW]) (*changesetmodels.ModelChanges, error) {
	recordChanges := changesetmodels.RecordChanges{}

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

func setModelIDOrCreate(modelChanges *changesetmodels.ModelChanges, schemaData *metadataclient.Schema, modelSpec spec.Model) error {
	if model, modelExists := schemaData.ModelByName(modelSpec.Name); modelExists {
		logger.Info("model exists", slog.String("modelName", modelSpec.Name), slog.String("modelID", model.ID))
		modelChanges.ID = model.ID
	} else {
		logger.Info("model must be created", slog.String("modelName", modelSpec.Name))
		propsCreate, err := modelSpec.PropertyCreator()
		if err != nil {
			return err
		}
		modelChanges.Create = &changesetmodels.ModelPropsCreate{
			Model: changesetmodels.ModelCreate{
				Name:        modelSpec.Name,
				DisplayName: modelSpec.DisplayName,
				Description: modelSpec.Description,
				Locked:      false,
			},
			Properties: propsCreate,
		}
	}
	return nil
}
