package models

type commonLabel struct {
	Id       string   `json:"id"`
	Label    string   `json:"label"`
	Synonyms []string `json:"synonyms,omitempty"`
	System   string   `json:"system"`
	Type     string   `json:"type"`
}

func newCommonLabel(id string, label string, system string, idType string, synonym ...string) commonLabel {
	return commonLabel{Id: id, Label: label, Synonyms: synonym, System: system, Type: idType}
}
