package fromcuration

import (
	"encoding/json"
	"fmt"
	"github.com/pennsieve/ttl-sync-processor/client/models/curation"
	"os"
)

func UnmarshalDatasetExport(datasetExportPath string) (*curation.DatasetExport, error) {
	curationExportFile, err := os.Open(datasetExportPath)
	if err != nil {
		return nil, fmt.Errorf("error opening curation export file %s: %w", datasetExportPath, err)
	}
	var curationExport curation.DatasetExport
	if err := json.NewDecoder(curationExportFile).Decode(&curationExport); err != nil {
		return nil, fmt.Errorf("error decoding curation export file %s: %w", datasetExportPath, err)
	}
	return &curationExport, nil
}
