package sync

import (
	"fmt"
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/spec"
	"log/slog"
)

func ComputeContributorsChanges(schemaData *metadataclient.Schema, old []metadata.Contributor, new []metadata.Contributor) (*changesetmodels.ModelChanges, error) {
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
	changes := &changesetmodels.ModelChanges{}
	if err := setModelIDOrCreate(changes, schemaData, spec.Contributor); err != nil {
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
		create := spec.ContributorInstance.Creator(contributor)
		changes.Records.Create = append(changes.Records.Create, create)
	}
	return changes, nil
}
