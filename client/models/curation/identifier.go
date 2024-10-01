package curation

type embeddedIdentifier struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	Type        string   `json:"type"`
	Synonyms    []string `json:"synonyms,omitempty"`
	System      string   `json:"system,omitempty"`
	Description string   `json:"description,omitempty"`
}

func newEmbeddedIdentifier(id, label, system, description string, synonym ...string) embeddedIdentifier {
	return embeddedIdentifier{ID: id, Label: label, Synonyms: synonym, System: system, Type: "identifier", Description: description}
}

func newDescriptionlessIdentifier(id string, label string, system string, synonym ...string) embeddedIdentifier {
	return newEmbeddedIdentifier(id, label, system, "", synonym...)
}
