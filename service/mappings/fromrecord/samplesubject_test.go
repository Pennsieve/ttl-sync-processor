package fromrecord

import (
	"fmt"
	"github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/metadatatest"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSampleStoreMapping(t *testing.T) {
	inputDirectory := "testdata/input"

	reader, err := client.NewReader(inputDirectory)
	require.NoError(t, err)

	idMap := NewRecordIDStore()
	samples, err := MapRecords(reader, metadata.SampleModelName, ToSample)
	require.NoError(t, err)
	for _, s := range samples {
		idMap.Add(metadata.SampleModelName, s.ExternalID(), s.GetPennsieveID())
	}
	subjects, err := MapRecords(reader, metadata.SubjectModelName, ToSubject)
	require.NoError(t, err)
	for _, s := range subjects {
		idMap.Add(metadata.SubjectModelName, s.ExternalID(), s.GetPennsieveID())
	}

	mapping := NewSampleStoreMapping(idMap)

	sampleSubjectLinks, err := MapLinkedProperties(reader, metadata.SampleSubjectLinkName, mapping)
	require.NoError(t, err)
	assert.Len(t, sampleSubjectLinks, 2)

	// Fist PennsieveLink
	{
		link := sampleSubjectLinks[0]
		metadatatest.AssertPennsieveInstanceIDEqual(t, "c148b5ae-10ff-4c41-bff6-0c0753a01e49", link.PennsieveID)
		metadatatest.AssertPennsieveInstanceIDEqual(t, "c148b5ae-10ff-4c41-bff6-0c0753a01e49", link.GetPennsieveID())

		metadatatest.AssertPennsieveInstanceIDEqual(t, "b66cbf32-cb9f-4126-8182-01bd00ad7b17", link.From)
		metadatatest.AssertPennsieveInstanceIDEqual(t, "c823ae0b-0c83-48be-9b0e-16690f6e675e", link.To)

		metadatatest.AssertExternalInstanceIDEqual(t, "967af4ee-eca9-4977-a74d-88713b82975f", link.SampleID)
		metadatatest.AssertExternalInstanceIDEqual(t, "a9ea0803-651b-4f0e-bda5-e2430e235e94", link.SubjectID)

		metadatatest.AssertExternalInstanceIDEqual(t, fmt.Sprintf("%s:%s", link.SampleID, link.SubjectID), link.ExternalID())
	}

	// Second PennsieveLink
	{
		link := sampleSubjectLinks[1]
		metadatatest.AssertPennsieveInstanceIDEqual(t, "88bbdc99-cfba-4bd1-800c-dcaa93742196", link.PennsieveID)
		metadatatest.AssertPennsieveInstanceIDEqual(t, "88bbdc99-cfba-4bd1-800c-dcaa93742196", link.GetPennsieveID())

		metadatatest.AssertPennsieveInstanceIDEqual(t, "60f21224-481c-4e29-a325-c896f184aebe", link.From)
		metadatatest.AssertPennsieveInstanceIDEqual(t, "b5ad14ab-f9e7-480b-b929-8e56db504181", link.To)

		metadatatest.AssertExternalInstanceIDEqual(t, "09d2a327-be38-403a-884d-a4d1d98b732c", link.SampleID)
		metadatatest.AssertExternalInstanceIDEqual(t, "f1027e6e-17cf-45d7-8b57-4c91bfd93fce", link.SubjectID)

		metadatatest.AssertExternalInstanceIDEqual(t, fmt.Sprintf("%s:%s", link.SampleID, link.SubjectID), link.ExternalID())
	}

}
