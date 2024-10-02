package fromrecord

// SchemaData is a map from model name to either:
// the model's ID string if the model exists in the dataset's metadata, or
// a pointer to changeset/models.ModelCreate if the model does not exist
type SchemaData map[string]any
