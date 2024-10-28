package sync

import (
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
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
	modelLogger.Info("change summary",
		slog.Int("createCount", len(modelChanges.Records.Create)),
		slog.Int("updateCount", len(modelChanges.Records.Update)),
		slog.Int("deleteCount", len(modelChanges.Records.Delete)),
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

	for _, toDelete := range oldToDelete {
		recordChanges.Delete = append(recordChanges.Delete, toDelete.GetPennsieveID())
	}
	if len(recordChanges.Create) == 0 && len(recordChanges.Delete) == 0 && len(recordChanges.Update) == 0 {
		return nil, nil
	}

	return &changesetmodels.ModelChanges{Records: recordChanges}, nil
}

func setModelIDOrCreate(modelChanges *changesetmodels.ModelChanges, schemaData *metadataclient.Schema, modelSpec spec.Model) error {
	if model, modelExists := schemaData.ModelByName(modelSpec.Name); modelExists {
		logger.Info("model exists", slog.String("modelName", modelSpec.Name), slog.String("modelID", model.ID))
		modelChanges.ID = changesetmodels.PennsieveSchemaID(model.ID)
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
