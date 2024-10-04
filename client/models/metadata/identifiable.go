package metadata

import changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"

type Saved interface {
	GetPennsieveID() changesetmodels.PennsieveRecordID
}

type IDer interface {
	GetID() string
}

type SavedIDer interface {
	Saved
	IDer
}
