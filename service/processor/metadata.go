package processor

import (
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/mappings/fromrecord"
)

func (p *CurationExportSyncProcessor) ExistingPennsieveMetadata(schemaData fromrecord.SchemaData) (*metadata.DatasetMetadata, error) {
	return fromrecord.ToDatasetMetadata(p.MetadataReader, schemaData)
}
