package spec

import "github.com/pennsieve/ttl-sync-processor/client/models/metadata"

var SampleSubject = Link{
	Name:          metadata.SampleSubjectLinkName,
	DisplayName:   metadata.SampleSubjectLinkDisplayName,
	FromModelName: metadata.SampleModelName,
	ToModelName:   metadata.SubjectModelName,
	Position:      1,
}
