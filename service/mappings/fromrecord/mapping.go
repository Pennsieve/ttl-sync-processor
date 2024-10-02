package fromrecord

import (
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
