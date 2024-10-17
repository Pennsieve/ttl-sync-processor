package spec

import (
	"fmt"
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-pre-metadata/client/models/datatypes"
)

type arrayPropertyCreator func(propertyName, displayName string, itemDataType datatypes.SimpleType) (changesetmodels.PropertyCreate, error)

func newArrayPropertyCreate(propertyName, displayName string, itemDataType datatypes.SimpleType) (changesetmodels.PropertyCreate, error) {
	propCreate := &changesetmodels.PropertyCreate{
		DisplayName: displayName,
		Name:        propertyName,
	}
	dataType := datatypes.ArrayDataType{
		Type: datatypes.ArrayType,
		Items: datatypes.ItemsType{
			Type: itemDataType,
		},
	}
	if err := propCreate.SetDataType(dataType); err != nil {
		return changesetmodels.PropertyCreate{}, fmt.Errorf("error setting data type of %s %s to %s: %w", propertyName,
			displayName, dataType, err)
	}
	return *propCreate, nil
}

func appendArrayPropertyCreate(creates []changesetmodels.PropertyCreate, propertyName, displayName string, itemDataType datatypes.SimpleType, propCreator arrayPropertyCreator, errs *[]error) []changesetmodels.PropertyCreate {
	create, err := propCreator(propertyName, displayName, itemDataType)
	if err != nil {
		*errs = append(*errs, err)
		return creates
	}
	creates = append(creates, create)
	return creates

}
