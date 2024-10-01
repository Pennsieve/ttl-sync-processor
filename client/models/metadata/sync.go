package metadata

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type Sync struct {
	Contributors []Contributor
}

func hashField(field any, fieldName string) (string, error) {
	bytes, err := json.Marshal(field)
	if err != nil {
		return "", fmt.Errorf("error marshalling %s: %w", fieldName, err)
	}
	hashBytes := md5.Sum(bytes)
	return hex.EncodeToString(hashBytes[:]), nil
}

func (s Sync) ContributorsHash() (string, error) {
	return hashField(s.Contributors, "contributors")
}
