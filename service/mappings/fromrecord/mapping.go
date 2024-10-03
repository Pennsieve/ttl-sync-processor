package fromrecord

import (
	"fmt"
	"github.com/pennsieve/processor-pre-metadata/client/models/instance"
	"github.com/pennsieve/ttl-sync-processor/service/logging"
)

var logger = logging.PackageLogger("fromrecord")

type Mapping[T any] func(record instance.Record) (T, error)

func MapSlice[T any](records []instance.Record, fromRecord Mapping[T]) ([]T, error) {
	var results []T
	for _, r := range records {
		result, err := fromRecord(r)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

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

func checkRecordType(record instance.Record, expectedModelName string) error {
	if record.Type != expectedModelName {
		return fmt.Errorf("record %s is not a %s instance: %s", record.ID, expectedModelName, record.Type)
	}
	return nil
}
