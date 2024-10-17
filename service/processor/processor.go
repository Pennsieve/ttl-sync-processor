package processor

import (
	"encoding/json"
	"fmt"
	changeset "github.com/pennsieve/processor-post-metadata/client"
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/service/logging"
	"github.com/pennsieve/ttl-sync-processor/service/sync"
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
	MetadataReader  *metadataclient.Reader
}

func NewCurationExportSyncProcessor(integrationID string, inputDirectory string, outputDirectory string) (*CurationExportSyncProcessor, error) {
	reader, err := metadataclient.NewReader(inputDirectory)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata reader for %s: %w", inputDirectory, err)
	}
	return &CurationExportSyncProcessor{
		IntegrationID:   integrationID,
		InputDirectory:  inputDirectory,
		OutputDirectory: outputDirectory,
		MetadataReader:  reader,
	}, nil
}

func (p *CurationExportSyncProcessor) Run() error {
	logger.Info("starting sync processing")
	logger.Info("Reading metadata from curation download")
	newMetadata, err := p.FromCuration()
	if err != nil {
		return err
	}
	logger.Info("Reading existing metadata from Pennsieve download")
	oldMetadata, err := p.ExistingPennsieveMetadata()
	if err != nil {
		return err
	}
	logger.Info("Computing required changes")
	changes, err := sync.ComputeChangeset(p.MetadataReader.Schema, oldMetadata, newMetadata)
	if err != nil {
		return err
	}
	if err := p.writeChangeset(changes); err != nil {
		return err
	}
	logger.Info("finished sync processing")
	return nil
}

func (p *CurationExportSyncProcessor) ChangesetFilePath() string {
	return filepath.Join(p.OutputDirectory, changeset.Filename)
}

func (p *CurationExportSyncProcessor) writeChangeset(changes *changesetmodels.Dataset) error {
	filePath := p.ChangesetFilePath()
	file, err := os.Create(filePath)
	defer func() {
		if err := file.Close(); err != nil {
			logger.Warn("error closing changeset file", slog.String("path", filePath), slog.Any("error", err))
		}
	}()
	if err != nil {
		return fmt.Errorf("error creating changeset file %s: %w", filePath, err)
	}
	if err := json.NewEncoder(file).Encode(changes); err != nil {
		return fmt.Errorf("error writing changeset file: %s: %w", filePath, err)
	}
	logger.Info("wrote changeset file", slog.String("path", filePath))
	return nil
}
