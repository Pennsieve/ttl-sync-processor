package sync

import (
	"github.com/pennsieve/processor-pre-metadata/client/models/schema"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/service/spec"
	"log/slog"
)

func setModelIDOrCreate(modelChanges *changesetmodels.ModelChanges, schemaData map[string]schema.Element, modelSpec spec.Model) error {
	if model, modelExists := schemaData[modelSpec.Name]; modelExists {
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
