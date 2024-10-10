package spec

import (
	"errors"
	"fmt"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
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
		return creates, errors.Join(errs...)
	},
}

var SampleInstance = IdentifiableInstance[metadata.SavedSample, metadata.Sample]{
	Creator: func(new metadata.Sample) changesetmodels.RecordCreate {
		var values []changesetmodels.RecordValue
		values = appendNonEmptyRecordValue(values, metadata.SampleIDKey, new.ExternalID())
		return changesetmodels.RecordCreate{Values: values}
	},
	Updater: func(old metadata.SavedSample, new metadata.Sample) (*changesetmodels.RecordUpdate, error) {
		if old.ExternalID() != new.ExternalID() {
			return nil, fmt.Errorf("old sample %s and new sample %s do not represent the same sample",
				old.ExternalID(),
				new.ExternalID())
		}
		// No updates, since this model does not have any non-ID properties at the moment. See subject for an example
		// of how this works with non-ID properties
		return nil, nil
	},
	Model: Sample,
}
