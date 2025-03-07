package sync

import (
	"fmt"
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/spec"
	"log/slog"
)

func ComputeContributorsChanges(schemaData *metadataclient.Schema, old []metadata.SavedContributor, new []metadata.Contributor) (any, error) {
	oldHash, err := metadata.ComputeHash(old)
	if err != nil {
		return nil, fmt.Errorf("error computing hash of existing contributors metadata: %w", err)
	}
	newHash, err := metadata.ComputeHash(new)
	if err != nil {
		return nil, fmt.Errorf("error computing hash of incoming contributors metadata: %w", err)
	}
	if oldHash == newHash {
		logger.Info("no changes required", slog.String("name", metadata.ContributorModelName))
		return nil, nil
	}
	// hashes are different, so clear out existing records and create new ones from incoming
	var toDelete []changesetmodels.PennsieveInstanceID
	for _, contributor := range old {
		toDelete = append(toDelete, contributor.GetPennsieveID())
	}

	var toCreate []changesetmodels.RecordCreate
	for _, contributor := range new {
		create := spec.ContributorInstance.Creator(contributor)
		toCreate = append(toCreate, create)
	}

	modelSpec := spec.Contributor
	modelName := modelSpec.Name
	modelLogger := logger.With(slog.String("modelName", modelName))
	if model, modelExists := schemaData.ModelByName(modelName); modelExists {
		modelLogger.With(
			slog.String("modelID", model.ID),
			slog.Int("toDeleteCount", len(toDelete)),
			slog.Int("toCreateCount", len(toCreate)),
		).Info("model exists")
		return &changesetmodels.ModelUpdate{
			ID: changesetmodels.PennsieveSchemaID(model.ID),
			Records: changesetmodels.RecordChanges{
				Delete: toDelete,
				Create: toCreate,
			},
		}, nil
	}

	modelLogger.With(slog.Int("toCreateCount", len(toCreate))).Info("model must be created")
	if len(toDelete) > 0 {
		return nil, fmt.Errorf("illegal state: a ModelCreate cannot contain record deletes: %s", modelName)
	}
	propsCreate, err := spec.Contributor.PropertyCreator()
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
		Records: toCreate,
	}, nil
}
