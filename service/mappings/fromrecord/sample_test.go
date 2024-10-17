package fromrecord

import (
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/metadatatest"
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
	assert.Len(t, samples, 3)

	sample1 := samples[0]
	metadatatest.AssertPennsieveInstanceIDEqual(t, "b66cbf32-cb9f-4126-8182-01bd00ad7b17", sample1.PennsieveID)
	metadatatest.AssertPennsieveInstanceIDEqual(t, "b66cbf32-cb9f-4126-8182-01bd00ad7b17", sample1.GetPennsieveID())

	metadatatest.AssertExternalInstanceIDEqual(t, "967af4ee-eca9-4977-a74d-88713b82975f", sample1.ID)
	metadatatest.AssertExternalInstanceIDEqual(t, "967af4ee-eca9-4977-a74d-88713b82975f", sample1.ExternalID())

	assert.Equal(t, "f61ed5a9-5a69-49f7-9113-6447ee9e668b", sample1.PrimaryKey)

	sample2 := samples[1]
	metadatatest.AssertPennsieveInstanceIDEqual(t, "60f21224-481c-4e29-a325-c896f184aebe", sample2.PennsieveID)
	metadatatest.AssertPennsieveInstanceIDEqual(t, "60f21224-481c-4e29-a325-c896f184aebe", sample2.GetPennsieveID())

	metadatatest.AssertExternalInstanceIDEqual(t, "09d2a327-be38-403a-884d-a4d1d98b732c", sample2.ID)
	metadatatest.AssertExternalInstanceIDEqual(t, "09d2a327-be38-403a-884d-a4d1d98b732c", sample2.ExternalID())

	assert.Equal(t, "rat-sample-pk-1", sample2.PrimaryKey)

	sample3 := samples[2]
	metadatatest.AssertPennsieveInstanceIDEqual(t, "cf811b54-bab7-49f1-b239-79c0f6cac29a", sample3.PennsieveID)
	metadatatest.AssertPennsieveInstanceIDEqual(t, "cf811b54-bab7-49f1-b239-79c0f6cac29a", sample3.GetPennsieveID())

	metadatatest.AssertExternalInstanceIDEqual(t, "sample-689", sample3.ID)
	metadatatest.AssertExternalInstanceIDEqual(t, "sample-689", sample3.ExternalID())

	assert.Equal(t, "dog-sample-pk-5", sample3.PrimaryKey)

}
