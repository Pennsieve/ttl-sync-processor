package processor

import (
	"encoding/json"
	"fmt"
	"github.com/pennsieve/ttl-sync-processor/client/models"
	"github.com/pennsieve/ttl-sync-processor/service/logging"
	"log/slog"
	"os"
	"path/filepath"
)

var logger = logging.PackageLogger("processor")

const CurationExportFilename = "curation-export.json"

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
	curationExportPath := p.curationExportPath()
	curationExport, err := p.readCurationExport()
	if err != nil {
		return err
	}
	logger.Info("read curation export file",
		slog.String("path", curationExportPath),
		slog.Any("datasetId", curationExport.ID),
	)
	logger.Debug("curation export file contents",
		slog.Any("curation-export", curationExport),
	)
	logger.Info("finished sync processing")
	return nil
}

func (p *CurationExportSyncProcessor) readCurationExport() (*models.DatasetCurationExport, error) {
	curationExportPath := p.curationExportPath()
	curationExportFile, err := os.Open(curationExportPath)
	if err != nil {
		return nil, fmt.Errorf("error opening curation export file %s: %w", curationExportPath, err)
	}
	var curationExport models.DatasetCurationExport
	if err := json.NewDecoder(curationExportFile).Decode(&curationExport); err != nil {
		return nil, fmt.Errorf("error decoding curation export file %s: %w", curationExportPath, err)
	}
	return &curationExport, nil
}

func (p *CurationExportSyncProcessor) curationExportPath() string {
	return filepath.Join(p.InputDirectory, CurationExportFilename)
}
