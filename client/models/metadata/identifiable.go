package metadata

import changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"

type Saved interface {
	GetPennsieveID() changesetmodels.PennsieveInstanceID
}

type IDer interface {
	GetID() string
}

type SavedIDer interface {
	Saved
	IDer
}
