package spec

import (
	"errors"
	"fmt"
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-pre-metadata/client/models/datatypes"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

var Subject = Model{
	Name:        metadata.SubjectModelName,
	DisplayName: metadata.SubjectDisplayName,
	Description: "Subjects in this dataset",
	PropertyCreator: func() (changesetmodels.PropertiesCreateParams, error) {
		var create []changesetmodels.PropertyCreateParams
		var accumulatedErrors []error
		create = appendSimplePropertyCreate(create, metadata.SpeciesKey, "Species", datatypes.StringType, newSimplePropertyCreate, &accumulatedErrors)
		create = appendConceptTitlePropertyCreate(create, metadata.SubjectIDKey, "ID", newConceptTitlePropertyCreate, &accumulatedErrors)
		create = appendArrayPropertyCreate(create, metadata.SpeciesSynonymsKey, "Species Synonyms", datatypes.StringType, newArrayPropertyCreate, &accumulatedErrors)
		return create, errors.Join(accumulatedErrors...)
	},
}

var SubjectInstance = IdentifiableInstance[metadata.SavedSubject, metadata.Subject]{
	Creator: func(subject metadata.Subject) changesetmodels.RecordCreate {
		var values []changesetmodels.RecordValue
		values = appendNonEmptyRecordValue(values, metadata.SubjectIDKey, subject.ID)
		values = appendNonEmptyRecordValue(values, metadata.SpeciesKey, subject.Species)
		values = appendNonEmptyRecordValue(values, metadata.SpeciesSynonymsKey, subject.SpeciesSynonyms)
		return changesetmodels.RecordCreate{
			ExternalID:   subject.ExternalID(),
			RecordValues: changesetmodels.RecordValues{Values: values},
		}
	},
	Updater: func(oldSubject metadata.SavedSubject, newSubject metadata.Subject) (*changesetmodels.RecordUpdate, error) {
		// since we are identifying old and new based on GetID(), it doesn't make sense to update the ID
		if oldSubject.ExternalID() != newSubject.ExternalID() {
			return nil, fmt.Errorf("old subject %s and new subject %s do not represent the same subject",
				oldSubject.ExternalID(),
				newSubject.ExternalID())
		}
		var values []changesetmodels.RecordValue
		noChange := true
		// The ID cannot be updated, but if there are other changes, we need to include all properties, even those
		// not changed
		values, noChange = appendExternalIDRecordValue(values, metadata.SubjectIDKey, oldSubject.ID, newSubject.ID, noChange)
		values, noChange = appendStringRecordValue(values, metadata.SpeciesKey, oldSubject.Species, newSubject.Species, noChange)
		values, noChange = appendStringSliceRecordValue(values, metadata.SpeciesSynonymsKey, oldSubject.SpeciesSynonyms, newSubject.SpeciesSynonyms, noChange)

		if !noChange {
			return &changesetmodels.RecordUpdate{PennsieveID: oldSubject.PennsieveID, RecordValues: changesetmodels.RecordValues{Values: values}}, nil
		}
		return nil, nil
	},
	Model: Subject,
}
