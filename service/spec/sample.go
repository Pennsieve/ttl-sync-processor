package spec

import (
	"errors"
	"fmt"
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-pre-metadata/client/models/datatypes"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

var Sample = Model{
	Name:        metadata.SampleModelName,
	DisplayName: metadata.SampleDisplayName,
	Description: "The samples in this dataset",
	PropertyCreator: func() (changesetmodels.PropertiesCreate, error) {
		var creates []changesetmodels.PropertyCreate
		var errs []error
		creates = appendConceptTitlePropertyCreate(creates, metadata.SampleIDKey, "ID", newConceptTitlePropertyCreate, &errs)
		creates = appendSimplePropertyCreate(creates, metadata.PrimaryKeyKey, "Primary Key", datatypes.StringType, newSimplePropertyCreate, &errs)
		return creates, errors.Join(errs...)
	},
}

var SampleInstance = IdentifiableInstance[metadata.SavedSample, metadata.Sample]{
	Creator: func(new metadata.Sample) changesetmodels.RecordCreate {
		var values []changesetmodels.RecordValue
		values = appendNonEmptyRecordValue(values, metadata.SampleIDKey, new.ExternalID())
		values = appendNonEmptyRecordValue(values, metadata.PrimaryKeyKey, new.PrimaryKey)
		return changesetmodels.RecordCreate{Values: values}
	},
	Updater: func(old metadata.SavedSample, new metadata.Sample) (*changesetmodels.RecordUpdate, error) {
		if old.ExternalID() != new.ExternalID() {
			return nil, fmt.Errorf("old sample %s and new sample %s do not represent the same sample",
				old.ExternalID(),
				new.ExternalID())
		}
		var values []changesetmodels.RecordValue
		noChange := true
		// The ID cannot be updated, but if there are other changes, we need to include all properties, even those
		// not changed
		values, noChange = appendExternalIDRecordValue(values, metadata.SampleIDKey, old.ID, new.ID, noChange)
		values, noChange = appendStringRecordValue(values, metadata.PrimaryKeyKey, old.PrimaryKey, new.PrimaryKey, noChange)

		if !noChange {
			return &changesetmodels.RecordUpdate{PennsieveID: old.PennsieveID, RecordValues: changesetmodels.RecordValues{Values: values}}, nil
		}
		return nil, nil
	},
	Model: Sample,
}
