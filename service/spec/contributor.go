package spec

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-pre-metadata/client/models/datatypes"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"strings"
)

var Contributor = Model{
	Name:        metadata.ContributorModelName,
	DisplayName: metadata.ContributorDisplayName,
	Description: "Contributors to this dataset",
	PropertyCreator: func() (changesetmodels.PropertiesCreateParams, error) {
		var create []changesetmodels.PropertyCreateParams
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

		// This is a hack to build an externalID for the contributor since it does not have a natural one
		stringValues := []string{contributor.FirstName, contributor.MiddleInitial, contributor.LastName, contributor.Degree, contributor.ORCID, contributor.NodeID}
		content := strings.Join(stringValues, ":")
		hashBytes := md5.Sum([]byte(content))
		externalID := hex.EncodeToString(hashBytes[:])

		create := changesetmodels.RecordCreate{
			ExternalID:   changesetmodels.ExternalInstanceID(externalID),
			RecordValues: changesetmodels.RecordValues{Values: values},
		}
		return create
	},
}
