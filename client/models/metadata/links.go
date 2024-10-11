package metadata

import changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"

type PennsieveLink interface {
	FromPennsieveID() changesetmodels.PennsieveInstanceID
	ToPennsieveID() changesetmodels.PennsieveInstanceID
}

type SavedLink interface {
	Saved
	PennsieveLink
}

type Link struct {
	From changesetmodels.PennsieveInstanceID
	To   changesetmodels.PennsieveInstanceID
}

func (l Link) FromPennsieveID() changesetmodels.PennsieveInstanceID {
	return l.From
}

func (l Link) ToPennsieveID() changesetmodels.PennsieveInstanceID {
	return l.To
}

type ExternalLink interface {
	ExternalIDer
	FromExternalID() changesetmodels.ExternalInstanceID
	ToExternalID() changesetmodels.ExternalInstanceID
}

type SavedExternalLink interface {
	SavedLink
	ExternalLink
}
