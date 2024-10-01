package processor

import (
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/logging"
	"log/slog"
	"path/filepath"
)

var logger = logging.PackageLogger("processor")

const CurationExportFilename = "curation-export.json"

const ContributorModelName = "contributor"

type CurationExportSyncProcessor struct {
	IntegrationID   string
	InputDirectory  string
	OutputDirectory string
}

func NewCurationExportSyncProcessor(integrationID string, inputDirectory string, outputDirectory string) *CurationExportSyncProcessor {
	return &CurationExportSyncProcessor{
		IntegrationID:   integrationID,
		InputDirectory:  inputDirectory,
		OutputDirectory: outputDirectory,
	}
}

func (p *CurationExportSyncProcessor) Run() error {
	logger.Info("starting sync processing")
	logger.Info("Reading metadata from curation download")
	_, err := p.FromCuration()
	if err != nil {
		return err
	}
	logger.Info("Reading existing metadata from Pennsieve download")
	_, err = p.ReadExistingPennsieveMetadata()
	if err != nil {
		return err
	}

	logger.Info("finished sync processing")
	return nil
}

func (p *CurationExportSyncProcessor) FromCuration() (*metadata.Sync, error) {
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
	return p.ConvertCurationExport(curationExport)
}

func (p *CurationExportSyncProcessor) curationExportPath() string {
	return filepath.Join(p.InputDirectory, CurationExportFilename)
}
