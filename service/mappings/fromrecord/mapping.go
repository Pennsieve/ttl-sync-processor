package fromrecord

import (
	"fmt"
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	"github.com/pennsieve/processor-pre-metadata/client/models/instance"
	"github.com/pennsieve/ttl-sync-processor/service/logging"
)

var logger = logging.PackageLogger("fromrecord")

func safeString(value any) string {
	if value == nil {
		return ""
	}
	return value.(string)
}

// If safeStringSlice is passed nil, it returns an empty []string
// If it is passed a []any, where the underlying elements are strings
// it returns a []string with those values.
// panics on anything else
func safeStringSlice(value any) []string {
	if value == nil {
		return []string{}
	}
	anySlice := value.([]any)
	strSlice := make([]string, len(anySlice))
	for i := range anySlice {
		strSlice[i] = anySlice[i].(string)
	}
	return strSlice
}

func safeExternalID(value any) changesetmodels.ExternalInstanceID {
	if value == nil {
		return changesetmodels.ExternalInstanceID("")
	}
	return changesetmodels.ExternalInstanceID(value.(string))
}

func checkRecordType(record instance.Record, expectedModelName string) error {
	if record.Type != expectedModelName {
		return fmt.Errorf("record %s is not a %s instance: %s", record.ID, expectedModelName, record.Type)
	}
	return nil
}

func checkLinkedPropertyName(linkedProperty instance.LinkedProperty, expectedLinkedPropertyName string) error {
	if linkedProperty.Name != expectedLinkedPropertyName {
		return fmt.Errorf("linked property %s is not named %s: %s", linkedProperty.ID, expectedLinkedPropertyName, linkedProperty.Name)
	}
	return nil
}
