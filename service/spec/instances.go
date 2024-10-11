package spec

import (
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

type Instance[OLD, NEW any] struct {
	Creator func(new NEW) changesetmodels.RecordCreate
	Updater func(old OLD, new NEW) (*changesetmodels.RecordUpdate, error)
}

type IdentifiableInstance[OLD metadata.SavedExternalIDer, NEW metadata.ExternalIDer] struct {
	Creator func(new NEW) changesetmodels.RecordCreate
	Updater func(old OLD, new NEW) (*changesetmodels.RecordUpdate, error)
	Model   Model
}
