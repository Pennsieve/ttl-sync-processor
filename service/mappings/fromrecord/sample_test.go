package fromrecord

import (
	"fmt"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	changesetmodels "github.com/pennsieve/ttl-sync-processor/client/changeset/models"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToSample(t *testing.T) {
	inputDirectory := "testdata/input"

	reader, err := metadataclient.NewReader(inputDirectory)
	require.NoError(t, err)

	samples, err := MapRecords(reader, metadata.SampleModelName, ToSample)
	require.NoError(t, err)
	assert.Len(t, samples, 2)

	sample1 := samples[0]
	assert.Equal(t, changesetmodels.PennsieveInstanceID("b66cbf32-cb9f-4126-8182-01bd00ad7b17"), sample1.PennsieveID)
	assert.Equal(t, changesetmodels.PennsieveInstanceID("b66cbf32-cb9f-4126-8182-01bd00ad7b17"), sample1.GetPennsieveID())

	assert.Equal(t, "967af4ee-eca9-4977-a74d-88713b82975f", sample1.ID)
	assert.Equal(t, "967af4ee-eca9-4977-a74d-88713b82975f", sample1.GetID())

	sample2 := samples[1]
	assert.Equal(t, changesetmodels.PennsieveInstanceID("60f21224-481c-4e29-a325-c896f184aebe"), sample2.PennsieveID)
	assert.Equal(t, changesetmodels.PennsieveInstanceID("60f21224-481c-4e29-a325-c896f184aebe"), sample2.GetPennsieveID())

	assert.Equal(t, "09d2a327-be38-403a-884d-a4d1d98b732c", sample2.ID)
	assert.Equal(t, "09d2a327-be38-403a-884d-a4d1d98b732c", sample2.GetID())

}

func TestNewSampleStoreMapping(t *testing.T) {
	inputDirectory := "testdata/input"

	reader, err := metadataclient.NewReader(inputDirectory)
	require.NoError(t, err)

	samples, err := MapRecords(reader, metadata.SampleModelName, ToSample)
	require.NoError(t, err)
	subjects, err := MapRecords(reader, metadata.SubjectModelName, ToSubject)
	require.NoError(t, err)

	mapping := NewSampleStoreMapping(samples, subjects)

	sampleSubjectLinks, err := MapLinkedProperties(reader, metadata.SampleSubjectLinkName, mapping)
	require.NoError(t, err)
	assert.Len(t, sampleSubjectLinks, 2)

	// Fist Link
	{
		link := sampleSubjectLinks[0]
		assert.Equal(t, changesetmodels.PennsieveInstanceID("c148b5ae-10ff-4c41-bff6-0c0753a01e49"), link.PennsieveID)
		assert.Equal(t, changesetmodels.PennsieveInstanceID("c148b5ae-10ff-4c41-bff6-0c0753a01e49"), link.GetPennsieveID())

		assert.Equal(t, changesetmodels.PennsieveInstanceID("b66cbf32-cb9f-4126-8182-01bd00ad7b17"), link.From)
		assert.Equal(t, changesetmodels.PennsieveInstanceID("c823ae0b-0c83-48be-9b0e-16690f6e675e"), link.To)

		assert.Equal(t, "967af4ee-eca9-4977-a74d-88713b82975f", link.SampleID)
		assert.Equal(t, "a9ea0803-651b-4f0e-bda5-e2430e235e94", link.SubjectID)

		assert.Equal(t, fmt.Sprintf("%s:%s", link.SampleID, link.SubjectID), link.GetID())
	}

	// Second Link
	{
		link := sampleSubjectLinks[1]
		assert.Equal(t, changesetmodels.PennsieveInstanceID("88bbdc99-cfba-4bd1-800c-dcaa93742196"), link.PennsieveID)
		assert.Equal(t, changesetmodels.PennsieveInstanceID("88bbdc99-cfba-4bd1-800c-dcaa93742196"), link.GetPennsieveID())

		assert.Equal(t, changesetmodels.PennsieveInstanceID("60f21224-481c-4e29-a325-c896f184aebe"), link.From)
		assert.Equal(t, changesetmodels.PennsieveInstanceID("b5ad14ab-f9e7-480b-b929-8e56db504181"), link.To)

		assert.Equal(t, "09d2a327-be38-403a-884d-a4d1d98b732c", link.SampleID)
		assert.Equal(t, "f1027e6e-17cf-45d7-8b57-4c91bfd93fce", link.SubjectID)

		assert.Equal(t, fmt.Sprintf("%s:%s", link.SampleID, link.SubjectID), link.GetID())
	}

}
