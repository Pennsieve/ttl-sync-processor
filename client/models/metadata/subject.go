package metadata

import changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"

const SubjectModelName = "subject"
const SubjectDisplayName = "Subject"

// Keys should match the json struct tag

const SubjectIDKey = "id"
const SpeciesKey = "species"
const SpeciesSynonymsKey = "species_synonyms"

type Subject struct {
	ID              string   `json:"id"`
	Species         string   `json:"species"`
	SpeciesSynonyms []string `json:"species_synonyms,omitempty"`
}

type SavedSubject struct {
	PennsieveID changesetmodels.PennsieveRecordID `json:"-"`
	Subject
}
