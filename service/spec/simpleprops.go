package spec

import (
	"fmt"
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-pre-metadata/client/models/datatypes"
)

type simplePropertyCreator func(propertyName, displayName string, dataType datatypes.SimpleType) (changesetmodels.PropertyCreateParams, error)
type conceptTitlePropertyCreator func(propertyName, displayName string) (changesetmodels.PropertyCreateParams, error)

func newSimplePropertyCreate(propertyName, displayName string, dataType datatypes.SimpleType) (changesetmodels.PropertyCreateParams, error) {
	propCreate := &changesetmodels.PropertyCreateParams{
		DisplayName: displayName,
		Name:        propertyName,
	}
	if err := propCreate.SetDataType(dataType); err != nil {
		return changesetmodels.PropertyCreateParams{}, fmt.Errorf("error setting data type of %s %s to %s: %w", propertyName,
			displayName, dataType, err)
	}
	return *propCreate, nil
}

// newConceptTitlePropertyCreate always creates a String property
func newConceptTitlePropertyCreate(propertyName, displayName string) (changesetmodels.PropertyCreateParams, error) {
	propCreate := &changesetmodels.PropertyCreateParams{
		DisplayName:  displayName,
		Name:         propertyName,
		ConceptTitle: true,
		Required:     true,
	}
	if err := propCreate.SetDataType(datatypes.StringType); err != nil {
		return changesetmodels.PropertyCreateParams{}, fmt.Errorf("error setting data type of %s %s to %s: %w", propertyName,
			displayName, datatypes.StringType, err)
	}
	return *propCreate, nil
}

func appendSimplePropertyCreate(creates []changesetmodels.PropertyCreateParams, propertyName, displayName string, dataType datatypes.SimpleType, propCreator simplePropertyCreator, errs *[]error) []changesetmodels.PropertyCreateParams {
	create, err := propCreator(propertyName, displayName, dataType)
	if err != nil {
		*errs = append(*errs, err)
		return creates
	}
	creates = append(creates, create)
	return creates

}

func appendConceptTitlePropertyCreate(creates []changesetmodels.PropertyCreateParams, propertyName, displayName string, propCreator conceptTitlePropertyCreator, errs *[]error) []changesetmodels.PropertyCreateParams {
	create, err := propCreator(propertyName, displayName)
	if err != nil {
		*errs = append(*errs, err)
		return creates
	}
	creates = append(creates, create)
	return creates

}
