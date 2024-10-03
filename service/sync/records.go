package sync

import changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"

// HasLen does not attempt to cover every possible type that the builtin len can be called on.
// Only needs to capture those types that we are using as data types for these particular models.
type HasLen interface {
	~string | ~[]string
}

// appendNonEmptyRecordValue only appends a new changesetmodels.RecordValue if len(value) > 0
func appendNonEmptyRecordValue[T HasLen](values []changesetmodels.RecordValue, name string, value T) []changesetmodels.RecordValue {
	if len(value) > 0 {
		return append(values, changesetmodels.RecordValue{
			Value: value,
			Name:  name,
		})
	}
	return values
}