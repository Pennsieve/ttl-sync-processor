package sync

import (
	"errors"
	"github.com/pennsieve/processor-pre-metadata/client/models/datatypes"
	"github.com/pennsieve/processor-pre-metadata/client/models/schema"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"log/slog"
)

func ComputeSubjectChanges(schemaData map[string]schema.Element, old []metadata.SavedSubject, new []metadata.Subject) (*changesetmodels.ModelChanges, error) {
	modelChanges, err := addChanges(old, new)
	if err != nil {
		return nil, err
	}
	modelLogger := logger.With(slog.String("modelName", metadata.SubjectModelName))
	if modelChanges == nil {
		modelLogger.Info("no changes")
		return nil, nil
	}
	if err := subjectSetIDOrCreate(modelChanges, schemaData); err != nil {
		return nil, err
	}
	deleteMessage := slog.Int("deleteCount", len(modelChanges.Records.Delete))
	if modelChanges.Records.DeleteAll {
		deleteMessage = slog.Int("deleteAllCount", len(old))
	}
	modelLogger.Info("change summary",
		deleteMessage,
		slog.Int("createCount", len(modelChanges.Records.Create)),
		slog.Int("updateCount", len(modelChanges.Records.Update)),
	)
	return modelChanges, nil
}

func addChanges(old []metadata.SavedSubject, new []metadata.Subject) (*changesetmodels.ModelChanges, error) {
	recordChanges := changesetmodels.RecordChanges{}

	oldByID := map[string]metadata.SavedSubject{}
	oldToDelete := map[string]metadata.SavedSubject{}
	for _, s := range old {
		oldByID[s.ID] = s
		oldToDelete[s.ID] = s
	}

	for _, newSubject := range new {
		newID := newSubject.ID
		if _, found := oldToDelete[newID]; found {
			delete(oldToDelete, newID)
		}
		if oldSubject, exists := oldByID[newID]; !exists {
			recordChanges.Create = append(recordChanges.Create, createSubjectRecord(newSubject))
		} else {
			recordUpdate, err := updateSubjectRecord(oldSubject, newSubject)
			if err != nil {
				return nil, err
			}
			if recordUpdate != nil {
				recordChanges.Update = append(recordChanges.Update, *recordUpdate)
			}
		}
	}

	if len(oldToDelete) > 0 {
		if len(oldToDelete) == len(old) {
			// use batch delete if we're going to delete all the existing records anyway
			recordChanges.DeleteAll = true
		} else {
			for _, toDelete := range oldToDelete {
				recordChanges.Delete = append(recordChanges.Delete, toDelete.PennsieveID)
			}
		}
	}
	if recordChanges.DeleteAll == false && len(recordChanges.Create) == 0 && len(recordChanges.Delete) == 0 && len(recordChanges.Update) == 0 {
		return nil, nil
	}

	return &changesetmodels.ModelChanges{Records: recordChanges}, nil
}

func createSubjectRecord(subject metadata.Subject) changesetmodels.RecordCreate {
	var values []changesetmodels.RecordValue
	values = appendNonEmptyRecordValue(values, metadata.SubjectIDKey, subject.ID)
	values = appendNonEmptyRecordValue(values, metadata.SpeciesKey, subject.Species)
	values = appendNonEmptyRecordValue(values, metadata.SpeciesSynonymsKey, subject.SpeciesSynonyms)
	return changesetmodels.RecordCreate{Values: values}
}

func updateSubjectRecord(oldSubject metadata.SavedSubject, newSubject metadata.Subject) (*changesetmodels.RecordUpdate, error) {
	var values []changesetmodels.RecordValue
	noChange := true
	values, noChange = appendStringRecordValue(values, metadata.SubjectIDKey, oldSubject.ID, newSubject.ID, noChange)
	values, noChange = appendStringRecordValue(values, metadata.SpeciesKey, oldSubject.Species, newSubject.Species, noChange)
	values, noChange = appendStringSliceRecordValue(values, metadata.SpeciesSynonymsKey, oldSubject.SpeciesSynonyms, newSubject.SpeciesSynonyms, noChange)

	if !noChange {
		return &changesetmodels.RecordUpdate{PennsieveID: oldSubject.PennsieveID, RecordValues: changesetmodels.RecordValues{Values: values}}, nil
	}
	return nil, nil
}

func subjectSetIDOrCreate(modelChanges *changesetmodels.ModelChanges, schemaData map[string]schema.Element) error {
	if model, modelExists := schemaData[metadata.SubjectModelName]; modelExists {
		logger.Info("model exists", slog.String("modelName", metadata.SubjectModelName), slog.String("modelID", model.ID))
		modelChanges.ID = model.ID
	} else {
		logger.Info("model must be created", slog.String("modelName", metadata.SubjectModelName))
		propsCreate, err := subjectPropertiesCreate()
		if err != nil {
			return err
		}
		modelChanges.Create = &changesetmodels.ModelPropsCreate{
			Model: changesetmodels.ModelCreate{
				Name:        metadata.SubjectModelName,
				DisplayName: metadata.SubjectDisplayName,
				Locked:      false,
			},
			Properties: propsCreate,
		}
	}
	return nil
}

func subjectPropertiesCreate() (changesetmodels.PropertiesCreate, error) {
	var create []changesetmodels.PropertyCreate
	var accumulatedErrors []error
	create = appendSimplePropertyCreate(create, metadata.SpeciesKey, "Species", datatypes.StringType, newSimplePropertyCreate, &accumulatedErrors)
	create = appendConceptTitlePropertyCreate(create, metadata.SubjectIDKey, "ID", newConceptTitlePropertyCreate, &accumulatedErrors)
	create = appendArrayPropertyCreate(create, metadata.SpeciesSynonymsKey, "Species Synonyms", datatypes.StringType, newArrayPropertyCreate, &accumulatedErrors)
	return create, errors.Join(accumulatedErrors...)
}
