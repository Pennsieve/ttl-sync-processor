package sync

import (
	"errors"
	"fmt"
	"github.com/pennsieve/processor-pre-metadata/client/models/datatypes"
	"github.com/pennsieve/processor-pre-metadata/client/models/schema"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"log/slog"
)

func ComputeContributorsChanges(schemaData map[string]schema.Element, old []metadata.Contributor, new []metadata.Contributor) (*changesetmodels.ModelChanges, error) {
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
	changes, err := contributorsModelChanges(schemaData)
	if err != nil {
		return nil, err
	}
	logger.Info("deleting and creating records",
		slog.String("ModelName", metadata.ContributorModelName),
		slog.Int("toDeleteCount", len(old)),
		slog.Int("toCreateCount", len(new)),
	)
	// hashes are different, so clear out existing records and create new ones from incoming
	changes.Records.DeleteAll = true
	for _, contributor := range new {
		create := contributorRecordCreate(contributor)
		changes.Records.Create = append(changes.Records.Create, create)
	}
	return changes, nil
}

func contributorsModelChanges(schemaData map[string]schema.Element) (*changesetmodels.ModelChanges, error) {
	changes := &changesetmodels.ModelChanges{}
	if model, modelExists := schemaData[metadata.ContributorModelName]; modelExists {
		logger.Info("model exists", slog.String("modelName", metadata.ContributorModelName), slog.String("modelID", model.ID))
		changes.ID = model.ID
	} else {
		logger.Info("model must be created", slog.String("modelName", metadata.ContributorModelName))
		propsCreate, err := contributorsPropertiesCreate()
		if err != nil {
			return nil, err
		}
		changes.Create = &changesetmodels.ModelPropsCreate{
			Model: changesetmodels.ModelCreate{
				Name:        metadata.ContributorModelName,
				DisplayName: metadata.ContributorDisplayName,
				Description: "Contributors to this dataset",
				Locked:      false,
			},
			Properties: propsCreate,
		}
	}
	return changes, nil
}

func contributorsPropertiesCreate() (changesetmodels.PropertiesCreate, error) {
	var create []changesetmodels.PropertyCreate
	var accumulatedErrors []error
	create = appendSimplePropertyCreate(create, metadata.FirstNameKey, "First Name", datatypes.StringType, newSimplePropertyCreate, &accumulatedErrors)
	create = appendSimplePropertyCreate(create, metadata.MiddleInitialKey, "Middle Initial", datatypes.StringType, newSimplePropertyCreate, &accumulatedErrors)
	create = appendConceptTitlePropertyCreate(create, metadata.LastNameKey, "Last Name", newConceptTitlePropertyCreate, &accumulatedErrors)
	create = appendSimplePropertyCreate(create, metadata.DegreeKey, "Degree", datatypes.StringType, newSimplePropertyCreate, &accumulatedErrors)
	create = appendSimplePropertyCreate(create, metadata.ORCIDKey, "ORCID", datatypes.StringType, newSimplePropertyCreate, &accumulatedErrors)
	create = appendSimplePropertyCreate(create, metadata.NodeIDKey, "Node ID", datatypes.StringType, newSimplePropertyCreate, &accumulatedErrors)
	return create, errors.Join(accumulatedErrors...)
}

func contributorRecordCreate(contributor metadata.Contributor) changesetmodels.RecordCreate {
	var values []changesetmodels.RecordValue
	values = appendNonEmptyRecordValue(values, metadata.FirstNameKey, contributor.FirstName)
	values = appendNonEmptyRecordValue(values, metadata.MiddleInitialKey, contributor.MiddleInitial)
	values = appendNonEmptyRecordValue(values, metadata.LastNameKey, contributor.LastName)
	values = appendNonEmptyRecordValue(values, metadata.DegreeKey, contributor.Degree)
	values = appendNonEmptyRecordValue(values, metadata.ORCIDKey, contributor.ORCID)
	values = appendNonEmptyRecordValue(values, metadata.NodeIDKey, contributor.NodeID)

	create := changesetmodels.RecordCreate{
		Values: values,
	}
	return create
}
