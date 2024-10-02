package sync

import (
	"errors"
	"fmt"
	"github.com/pennsieve/processor-pre-metadata/client/models/schema"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"log/slog"
)

func ComputeContributorsChanges(schemaData map[string]schema.Element, old, new *metadata.DatasetMetadata) (*changesetmodels.ModelChanges, error) {
	oldHash, err := metadata.ComputeHash(old.Contributors)
	if err != nil {
		return nil, fmt.Errorf("error computing hash of existing contributors metadata: %w", err)
	}
	newHash, err := metadata.ComputeHash(new.Contributors)
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
		slog.String("name", metadata.ContributorModelName),
		slog.Int("toDeleteCount", len(old.Contributors)),
		slog.Int("toCreateCount", len(new.Contributors)),
	)
	// hashes are different, so clear out existing records and create new ones from incoming
	changes.Records.DeleteAll = true
	for _, contributor := range new.Contributors {
		create := contributorRecordCreate(contributor)
		changes.Records.Create = append(changes.Records.Create, create)
	}
	return changes, nil
}

func contributorsModelChanges(schemaData map[string]schema.Element) (*changesetmodels.ModelChanges, error) {
	changes := &changesetmodels.ModelChanges{}
	if model, modelExists := schemaData[metadata.ContributorModelName]; modelExists {
		logger.Info("model exists", slog.String("name", metadata.ContributorModelName), slog.String("id", model.ID))
		changes.ID = model.ID
	} else {
		logger.Info("model must be created", slog.String("name", metadata.ContributorModelName))
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
	create = appendStringPropertyCreate(create, metadata.FirstNameKey, "First Name", newStringPropertyCreate, &accumulatedErrors)
	create = appendStringPropertyCreate(create, metadata.MiddleInitialKey, "Middle Initial", newStringPropertyCreate, &accumulatedErrors)
	create = appendStringPropertyCreate(create, metadata.LastNameKey, "Last Name", newStringConceptTitlePropertyCreate, &accumulatedErrors)
	create = appendStringPropertyCreate(create, metadata.DegreeKey, "Degree", newStringPropertyCreate, &accumulatedErrors)
	create = appendStringPropertyCreate(create, metadata.ORCIDKey, "ORCID", newStringPropertyCreate, &accumulatedErrors)
	create = appendStringPropertyCreate(create, metadata.NodeIDKey, "Node ID", newStringPropertyCreate, &accumulatedErrors)
	return create, errors.Join(accumulatedErrors...)
}

func contributorRecordCreate(contributor metadata.Contributor) changesetmodels.RecordCreate {
	var values []changesetmodels.RecordValue
	values = appendStringRecordValue(values, metadata.FirstNameKey, contributor.FirstName)
	values = appendStringRecordValue(values, metadata.MiddleInitialKey, contributor.MiddleInitial)
	values = appendStringRecordValue(values, metadata.LastNameKey, contributor.LastName)
	values = appendStringRecordValue(values, metadata.DegreeKey, contributor.Degree)
	values = appendStringRecordValue(values, metadata.ORCIDKey, contributor.ORCID)
	values = appendStringRecordValue(values, metadata.NodeIDKey, contributor.NodeID)

	create := changesetmodels.RecordCreate{
		Values: values,
	}
	return create
}
