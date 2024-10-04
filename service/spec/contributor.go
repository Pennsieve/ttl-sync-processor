package spec

import (
	"errors"
	"github.com/pennsieve/processor-pre-metadata/client/models/datatypes"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

var Contributor = Model{
	Name:        metadata.ContributorModelName,
	DisplayName: metadata.ContributorDisplayName,
	Description: "Contributors to this dataset",
	PropertyCreator: func() (changesetmodels.PropertiesCreate, error) {
		var create []changesetmodels.PropertyCreate
		var accumulatedErrors []error
		create = appendSimplePropertyCreate(create, metadata.FirstNameKey, "First Name", datatypes.StringType, newSimplePropertyCreate, &accumulatedErrors)
		create = appendSimplePropertyCreate(create, metadata.MiddleInitialKey, "Middle Initial", datatypes.StringType, newSimplePropertyCreate, &accumulatedErrors)
		create = appendConceptTitlePropertyCreate(create, metadata.LastNameKey, "Last Name", newConceptTitlePropertyCreate, &accumulatedErrors)
		create = appendSimplePropertyCreate(create, metadata.DegreeKey, "Degree", datatypes.StringType, newSimplePropertyCreate, &accumulatedErrors)
		create = appendSimplePropertyCreate(create, metadata.ORCIDKey, "ORCID", datatypes.StringType, newSimplePropertyCreate, &accumulatedErrors)
		create = appendSimplePropertyCreate(create, metadata.NodeIDKey, "Node ID", datatypes.StringType, newSimplePropertyCreate, &accumulatedErrors)
		return create, errors.Join(accumulatedErrors...)
	},
}

var ContributorInstance = Instance[metadata.Contributor, metadata.Contributor]{
	Updater: nil,
	Creator: func(contributor metadata.Contributor) changesetmodels.RecordCreate {
		var values []changesetmodels.RecordValue
		values = appendNonEmptyRecordValue(values, metadata.FirstNameKey, contributor.FirstName)
		values = appendNonEmptyRecordValue(values, metadata.MiddleInitialKey, contributor.MiddleInitial)
		values = appendNonEmptyRecordValue(values, metadata.LastNameKey, contributor.LastName)
		values = appendNonEmptyRecordValue(values, metadata.DegreeKey, contributor.Degree)
		values = appendNonEmptyRecordValue(values, metadata.ORCIDKey, contributor.ORCID)
		values = appendNonEmptyRecordValue(values, metadata.NodeIDKey, contributor.NodeID)

		create := changesetmodels.RecordCreate{
			Values: values,
		}
		return create
	},
}
