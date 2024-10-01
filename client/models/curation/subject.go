package curation

import (
	"encoding/json"
	"errors"
	"fmt"
)

type SpeciesIdentifier embeddedIdentifier

type Subject struct {
	ID string `json:"subject_id"`
	// Species should unmarshal to either a string or SpeciesIdentifier
	// Use GetSpecies() to access value
	Species json.RawMessage `json:"species,omitempty"`
}

// GetSpecies either returns a string or a SpeciesIdentifier. The SpeciesIdentifier will be emtpy if the Subject had no species
// or if the Species value did not match either format
func (s Subject) GetSpecies() (any, error) {
	if s.Species == nil {
		return "", nil
	}
	var speciesString string
	if strErr := json.Unmarshal(s.Species, &speciesString); strErr != nil {
		var speciesIdentifier SpeciesIdentifier
		if idErr := json.Unmarshal(s.Species, &speciesIdentifier); idErr != nil {
			return nil, fmt.Errorf("error unmarshalling species value %s: %w", s.Species, errors.Join(strErr, idErr))
		}
		return speciesIdentifier, nil
	}
	return speciesString, nil
}
