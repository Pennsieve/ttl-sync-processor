package fromrecord

import (
	"github.com/pennsieve/processor-pre-metadata/client/models/instance"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

func ToContributor(record instance.Record) (metadata.Contributor, error) {
	contributor := metadata.Contributor{}
	if err := checkRecordType(record, metadata.ContributorModelName); err != nil {
		return contributor, err
	}
	for _, v := range record.Values {
		switch v.Name {
		case metadata.FirstNameKey:
			contributor.FirstName = safeString(v.Value)
		case metadata.LastNameKey:
			contributor.LastName = safeString(v.Value)
		case metadata.MiddleInitialKey:
			contributor.MiddleInitial = safeString(v.Value)
		case metadata.DegreeKey:
			contributor.Degree = safeString(v.Value)
		case metadata.NodeIDKey:
			contributor.NodeID = safeString(v.Value)
		case metadata.ORCIDKey:
			contributor.ORCID = safeString(v.Value)
		}
	}
	return contributor, nil
}
