package metadata

import changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"

type Saved interface {
	GetPennsieveID() changesetmodels.PennsieveInstanceID
}

type ExternalIDer interface {
	ExternalID() changesetmodels.ExternalInstanceID
}

type SavedExternalIDer interface {
	Saved
	ExternalIDer
}
