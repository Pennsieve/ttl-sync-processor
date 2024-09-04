package modelstest

import "github.com/pennsieve/ttl-sync-processor/client/models"

type DatasetCurationExportBuilder struct {
	curationExport *models.DatasetCurationExport
}

func NewDatasetCurationExportBuilder(id string) *DatasetCurationExportBuilder {
	return &DatasetCurationExportBuilder{curationExport: &models.DatasetCurationExport{ID: id}}
}

func (b *DatasetCurationExportBuilder) WithContributors(contributors ...models.Contributor) *DatasetCurationExportBuilder {
	b.curationExport.Contributors = append(b.curationExport.Contributors, contributors...)
	return b
}

func (b *DatasetCurationExportBuilder) Build() models.DatasetCurationExport {
	return *b.curationExport
}
