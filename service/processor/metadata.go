package processor

import (
	"fmt"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/processor-pre-metadata/client/models/instance"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

func (p *CurationExportSyncProcessor) ReadExistingPennsieveMetadata() (*metadata.Sync, error) {
	reader, err := metadataclient.NewReader(p.InputDirectory)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata reader: %w", err)
	}
	existing := &metadata.Sync{}
	contributorRecords, err := reader.GetRecordsForModel(ContributorModelName)
	if err != nil {
		return nil, fmt.Errorf("error reading contributor records: %w", err)
	}
	existing.Contributors, err = FromRecords(contributorRecords, Contributor)
	if err != nil {
		return nil, fmt.Errorf("error marshalling contributors: %w", err)
	}

	return existing, nil
}

type FromRecord[T any] func(record instance.Record) (T, error)

func Contributor(record instance.Record) (metadata.Contributor, error) {
	contributor := metadata.Contributor{}
	if record.Type != ContributorModelName {
		return contributor, fmt.Errorf("record %s is not a contributor instance: %s", record.ID, record.Type)
	}
	for _, v := range record.Values {
		switch v.Name {
		case metadata.FirstNameKey:
			contributor.FirstName = safeString(v.Value)
		case metadata.LastNameKey:
			contributor.LastName = safeString(v.Value)
		case metadata.MiddleInitialKey:
			contributor.MiddleInitial = safeString(v.Value)
		case metadata.DegreeKey:
			contributor.Degree = safeString(v.Value)
		case metadata.NodeIDKey:
			contributor.NodeID = safeString(v.Value)
		case metadata.ORCIDKey:
			contributor.ORCID = safeString(v.Value)
		}
	}
	return contributor, nil
}

func FromRecords[T any](records []instance.Record, fromRecord FromRecord[T]) ([]T, error) {
	var results []T
	for _, r := range records {
		result, err := fromRecord(r)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func safeString(value any) string {
	if value == nil {
		return ""
	}
	return value.(string)
}
