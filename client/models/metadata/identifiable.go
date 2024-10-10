package metadata

import changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"

type Saved interface {
	GetPennsieveID() changesetmodels.PennsieveInstanceID
}

type ExternalIDer interface {
	ExternalID() changesetmodels.ExternalInstanceID
}

type SavedIDer interface {
	Saved
	ExternalIDer
}
