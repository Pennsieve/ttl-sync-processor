package metadata

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type DatasetMetadata struct {
	Contributors   []Contributor
	Subjects       []Subject
	Samples        []Sample
	SampleSubjects []SampleSubject
	Proxies        []Proxy
}

func computeHash(value any) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	hashBytes := md5.Sum(bytes)
	return hex.EncodeToString(hashBytes[:]), nil
}
func ComputeHash(value any) (string, error) {
	hash, err := computeHash(value)
	if err != nil {
		return "", fmt.Errorf("error marshalling value for hashing: %w", err)
	}
	return hash, nil
}

func hashField(field any, fieldName string) (string, error) {
	hash, err := computeHash(field)
	if err != nil {
		return "", fmt.Errorf("error marshalling %s: %w", fieldName, err)
	}
	return hash, nil
}

func (s DatasetMetadata) ContributorsHash() (string, error) {
	return hashField(s.Contributors, "contributors")
}

type SavedDatasetMetadata struct {
	Contributors   []Contributor
	Subjects       []SavedSubject
	Samples        []SavedSample
	SampleSubjects []SavedSampleSubject
	Proxies        []SavedProxy
}
