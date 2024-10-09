package sync

import (
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/service/spec"
	"log/slog"
)

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
