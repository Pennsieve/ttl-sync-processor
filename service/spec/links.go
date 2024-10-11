package spec

import changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"

type Link struct {
	Name          string
	DisplayName   string
	FromModelName string
	ToModelName   string
	Position      int
	SchemaCreator func() changesetmodels.SchemaLinkedPropertyCreate
}

func (l Link) SchemaCreate() changesetmodels.SchemaLinkedPropertyCreate {
	return changesetmodels.SchemaLinkedPropertyCreate{
		Name:          l.Name,
		DisplayName:   l.DisplayName,
		FromModelName: l.FromModelName,
		ToModelName:   l.ToModelName,
		Position:      l.Position,
	}
}
