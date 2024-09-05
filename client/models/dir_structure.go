package models

// DirStructureEntry represents an element in the dir_structure array.
// In the curation-export.json file it has many more properties and most
// are being ignored here
type DirStructureEntry struct {
	DatasetRelativePath string `json:"dataset_relative_path"`
	RemoteID            string `json:"remote_id"`
}

func NewDirStructureEntry(datasetRelativePath string, remoteID string) DirStructureEntry {
	return DirStructureEntry{
		DatasetRelativePath: datasetRelativePath,
		RemoteID:            remoteID,
	}
}
