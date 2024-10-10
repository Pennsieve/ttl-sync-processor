package metadata

import changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"

const SubjectModelName = "subject"
const SubjectDisplayName = "Subject"

// Keys should match the json struct tag

const SubjectIDKey = "id"
const SpeciesKey = "species"
const SpeciesSynonymsKey = "species_synonyms"

type Subject struct {
	ID              changesetmodels.ExternalInstanceID `json:"id"`
	Species         string                             `json:"species"`
	SpeciesSynonyms []string                           `json:"species_synonyms,omitempty"`
}

func (s Subject) ExternalID() changesetmodels.ExternalInstanceID {
	return changesetmodels.ExternalInstanceID(s.ID)
}

type SavedSubject struct {
	PennsieveID changesetmodels.PennsieveInstanceID `json:"-"`
	Subject
}

func (ss SavedSubject) GetPennsieveID() changesetmodels.PennsieveInstanceID {
	return ss.PennsieveID
}
