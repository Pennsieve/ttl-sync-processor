package curation

// DirStructureEntry represents an element in the dir_structure array.
// These are the directories in the dataset.
// In the curation-export.json file it has many more properties and most
// are being ignored here
type DirStructureEntry struct {
	// DatasetRelativePath is the full path of this directory entry
	DatasetRelativePath string `json:"dataset_relative_path"`
	// RemoteID is the Pennsieve package node id of this directory entry
	RemoteID string `json:"remote_id"`
}

func NewDirStructureEntry(datasetRelativePath string, remoteID string) DirStructureEntry {
	return DirStructureEntry{
		DatasetRelativePath: datasetRelativePath,
		RemoteID:            remoteID,
	}
}
