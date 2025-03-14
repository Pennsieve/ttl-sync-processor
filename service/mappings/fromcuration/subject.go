package fromcuration

import (
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/curation"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
)

// ToSubject is a mappings.Mapping from curation.Subject to metadata.Subject
func ToSubject(exportedSubject curation.Subject) (metadata.Subject, error) {
	speciesRaw, err := exportedSubject.GetSpecies()
	if err != nil {
		return metadata.Subject{}, err
	}
	var speciesName string
	var synonyms []string
	switch species := speciesRaw.(type) {
	case string:
		speciesName = species
	case curation.SpeciesIdentifier:
		speciesName = species.Label
		synonyms = species.Synonyms
	default:
		speciesName = "unknown"
	}
	if len(speciesName) == 0 {
		speciesName = "unknown"
	}
	subject := metadata.Subject{
		ID:              changesetmodels.ExternalInstanceID(exportedSubject.ID),
		Species:         speciesName,
		SpeciesSynonyms: synonyms,
	}
	return subject, nil
}
