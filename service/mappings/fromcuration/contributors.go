package fromcuration

import (
	"fmt"
	"github.com/pennsieve/ttl-sync-processor/client/models/curation"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"strings"
	"unicode/utf8"
)

// ToContributor is a Mapping from curation.Contributor to metadata.Contributor
func ToContributor(exportedContributor curation.Contributor) (metadata.Contributor, error) {
	contrib := metadata.Contributor{
		FirstName: exportedContributor.FirstName,
		LastName:  exportedContributor.LastName,
	}
	if len(exportedContributor.MiddleName) > 0 {
		initialRune, sz := utf8.DecodeRuneInString(exportedContributor.MiddleName)
		if initialRune == utf8.RuneError && sz == 0 {
			return metadata.Contributor{}, fmt.Errorf("impossible error for non-empty middle name")
		} else if initialRune == utf8.RuneError && sz == 1 {
			return metadata.Contributor{}, fmt.Errorf("middle name %s is not utf-8 encoded", exportedContributor.MiddleName)
		}
		contrib.MiddleInitial = fmt.Sprintf("%c", initialRune)
	}
	if exportedContributor.ContributorORCID != nil {
		contrib.ORCID = exportedContributor.ContributorORCID.ID
	}
	if strings.HasPrefix(exportedContributor.DataRemoteUserID, "user:") {
		contrib.NodeID = fmt.Sprintf("N:%s", exportedContributor.DataRemoteUserID)
	}
	return contrib, nil
}
