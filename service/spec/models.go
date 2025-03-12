package spec

import changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"

// Model contains information needed to create a model in the metadata schema
type Model struct {
	Name            string
	DisplayName     string
	Description     string
	PropertyCreator func() (changesetmodels.PropertiesCreateParams, error)
}
