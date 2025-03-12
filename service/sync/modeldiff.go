package sync

import (
	"fmt"
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
	instanceSpec spec.IdentifiableInstance[OLD, NEW]) (any, error) {
	recordChanges, err := addIdentifiableModelChanges(old, new, instanceSpec)
	if err != nil {
		return nil, err
	}
	modelLogger := logger.With(slog.String("modelName", instanceSpec.Model.Name))
	if recordChanges == nil {
		modelLogger.Info("no changes")
		return nil, nil
	}
	modelChanges, err := setModelIDOrCreate(*recordChanges, schemaData, instanceSpec.Model)
	if err != nil {
		return nil, err
	}
	modelLogger.Info("change summary",
		slog.Int("createCount", len(recordChanges.Create)),
		slog.Int("updateCount", len(recordChanges.Update)),
		slog.Int("deleteCount", len(recordChanges.Delete)),
	)
	return modelChanges, nil
}
func addIdentifiableModelChanges[OLD metadata.SavedExternalIDer, NEW metadata.ExternalIDer](old []OLD, new []NEW, instanceSpec spec.IdentifiableInstance[OLD, NEW]) (*changesetmodels.RecordChanges, error) {
	recordChanges := &changesetmodels.RecordChanges{}

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

	return recordChanges, nil
}

func setModelIDOrCreate(recordChanges changesetmodels.RecordChanges, schemaData *metadataclient.Schema, modelSpec spec.Model) (any, error) {
	if model, modelExists := schemaData.ModelByName(modelSpec.Name); modelExists {
		logger.Info("model exists", slog.String("modelName", modelSpec.Name), slog.String("modelID", model.ID))
		return &changesetmodels.ModelUpdate{
			ID:      changesetmodels.PennsieveSchemaID(model.ID),
			Records: recordChanges,
		}, nil
	}
	logger.Info("model must be created", slog.String("modelName", modelSpec.Name))
	if len(recordChanges.Delete) > 0 || len(recordChanges.Update) > 0 {
		return nil, fmt.Errorf("illegal state: a ModelCreate cannot contain record deletes or updates: %s", modelSpec.Name)
	}
	propsCreate, err := modelSpec.PropertyCreator()
	if err != nil {
		return nil, err
	}
	return &changesetmodels.ModelCreate{
		Create: changesetmodels.ModelPropsCreate{
			Model: changesetmodels.ModelCreateParams{
				Name:        modelSpec.Name,
				DisplayName: modelSpec.DisplayName,
				Description: modelSpec.Description,
				Locked:      false,
			},
			Properties: propsCreate,
		},
		Records: recordChanges.Create,
	}, nil

}
