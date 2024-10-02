package processor

import (
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/mappings/fromrecord"
)

func (p *CurationExportSyncProcessor) ExistingPennsieveMetadata() (*metadata.DatasetMetadata, error) {
	return fromrecord.ToDatasetMetadata(p.MetadataReader)
}
