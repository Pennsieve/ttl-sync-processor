package processor

import (
	"github.com/pennsieve/ttl-sync-processor/client/models/curation"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/mappings/fromcuration"
	"log/slog"
	"path/filepath"
)

func (p *CurationExportSyncProcessor) FromCuration() (*metadata.DatasetMetadata, error) {
	curationExportPath := p.curationExportPath()
	curationExport, err := p.ReadCurationExport()
	if err != nil {
		return nil, err
	}
	logger.Info("read curation export file",
		slog.String("path", curationExportPath),
		slog.Any("datasetId", curationExport.ID),
	)
	logger.Debug("curation export file contents",
		slog.Any("curation-export", curationExport),
	)
	return fromcuration.ToDatasetMetadata(curationExport)
}

func (p *CurationExportSyncProcessor) ReadCurationExport() (*curation.DatasetExport, error) {
	return fromcuration.UnmarshalDatasetExport(p.curationExportPath())
}

func (p *CurationExportSyncProcessor) curationExportPath() string {
	return filepath.Join(p.InputDirectory, CurationExportFilename)
}
