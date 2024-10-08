package metadata

import changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"

type link struct {
	From changesetmodels.PennsieveInstanceID
	To   changesetmodels.PennsieveInstanceID
}
