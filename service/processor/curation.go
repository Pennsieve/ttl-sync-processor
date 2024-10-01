package processor

import (
	"encoding/json"
	"fmt"
	"github.com/pennsieve/ttl-sync-processor/client/models/curation"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"os"
	"unicode/utf8"
)

func (p *CurationExportSyncProcessor) ReadCurationExport() (*curation.DatasetExport, error) {
	curationExportPath := p.curationExportPath()
	curationExportFile, err := os.Open(curationExportPath)
	if err != nil {
		return nil, fmt.Errorf("error opening curation export file %s: %w", curationExportPath, err)
	}
	var curationExport curation.DatasetExport
	if err := json.NewDecoder(curationExportFile).Decode(&curationExport); err != nil {
		return nil, fmt.Errorf("error decoding curation export file %s: %w", curationExportPath, err)
	}
	return &curationExport, nil
}

func (p *CurationExportSyncProcessor) ConvertCurationExport(export *curation.DatasetExport) (*metadata.Sync, error) {
	exported := &metadata.Sync{}
	var err error
	exported.Contributors, err = FromCurationExport(export.Contributors, ContributorFromCurationExport)
	if err != nil {
		return nil, err
	}
	return exported, nil
}

type CurationMapping[From, To any] func(curationObject From) (To, error)

func ContributorFromCurationExport(exportedContributor curation.Contributor) (metadata.Contributor, error) {
	contrib := metadata.Contributor{
		FirstName: exportedContributor.FirstName,
		LastName:  exportedContributor.LastName,
	}
	if len(exportedContributor.MiddleName) > 0 {
		initialRune, sz := utf8.DecodeRuneInString(exportedContributor.MiddleName)
		if initialRune == utf8.RuneError && sz == 0 {
			return metadata.Contributor{}, fmt.Errorf("impossible error for non-empty middle name")
		} else if initialRune == utf8.RuneError && sz == 1 {
			return metadata.Contributor{}, fmt.Errorf("middle name %s is not utf-8 encoded", exportedContributor.MiddleName)
		}
		contrib.MiddleInitial = fmt.Sprintf("%c", initialRune)
	}
	return contrib, nil
}

func FromCurationExport[From, To any](exported []From, mapping CurationMapping[From, To]) ([]To, error) {
	var pennsieveMetadata []To
	for _, e := range exported {
		m, err := mapping(e)
		if err != nil {
			return nil, err
		}
		pennsieveMetadata = append(pennsieveMetadata, m)
	}
	return pennsieveMetadata, nil
}
