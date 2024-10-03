package sync

import (
	"github.com/pennsieve/processor-pre-metadata/client/models/schema"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

func ComputeSubjectChanges(schemaData map[string]schema.Element, old, new *metadata.DatasetMetadata) (*changesetmodels.ModelChanges, error) {
	modelChanges, err := addChanges(old.Subjects, new.Subjects)
	if err != nil {
		return nil, err
	}
	if modelChanges == nil {
		return nil, nil
	}
	return modelChanges, nil
}

func addChanges(old, new []metadata.Subject) (*changesetmodels.ModelChanges, error) {
	recordChanges := changesetmodels.RecordChanges{}

	oldByID := map[string]metadata.Subject{}
	oldToDelete := map[string]bool{}
	for _, s := range old {
		oldByID[s.ID] = s
		oldToDelete[s.ID] = true
	}

	toDeleteCount := len(old)
	for _, newSubject := range new {
		newID := newSubject.ID
		if _, found := oldToDelete[newID]; found {
			oldToDelete[newID] = false
			toDeleteCount--
		}
		if oldSubject, exists := oldByID[newID]; !exists {
			recordCreate, err := createSubjectRecord(newSubject)
			if err != nil {
				return nil, err
			}
			recordChanges.Create = append(recordChanges.Create, recordCreate)
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

	if toDeleteCount > 0 {
		if toDeleteCount == len(old) {
			// use batch delete if we're going to delete all the existing records anyway
			recordChanges.DeleteAll = true
		} else {
			for id, doDelete := range oldToDelete {
				if doDelete {
					recordChanges.Delete = append(recordChanges.Delete, id)
				}
			}
		}
	}
	if recordChanges.DeleteAll == false && len(recordChanges.Create) == 0 && len(recordChanges.Delete) == 0 && len(recordChanges.Update) == 0 {
		return nil, nil
	}

	return &changesetmodels.ModelChanges{Records: recordChanges}, nil
}

func createSubjectRecord(subject metadata.Subject) (changesetmodels.RecordCreate, error) {
	create := changesetmodels.RecordCreate{}
	return create, nil
}

func updateSubjectRecord(oldSubject, newSubject metadata.Subject) (*changesetmodels.RecordUpdate, error) {
	update := &changesetmodels.RecordUpdate{}
	return update, nil
}
