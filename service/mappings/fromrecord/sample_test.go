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
	assert.Len(t, samples, 2)

	sample1 := samples[0]
	metadatatest.AssertPennsieveInstanceIDEqual(t, "b66cbf32-cb9f-4126-8182-01bd00ad7b17", sample1.PennsieveID)
	metadatatest.AssertPennsieveInstanceIDEqual(t, "b66cbf32-cb9f-4126-8182-01bd00ad7b17", sample1.GetPennsieveID())

	metadatatest.AssertExternalInstanceIDEqual(t, "967af4ee-eca9-4977-a74d-88713b82975f", sample1.ID)
	metadatatest.AssertExternalInstanceIDEqual(t, "967af4ee-eca9-4977-a74d-88713b82975f", sample1.ExternalID())

	sample2 := samples[1]
	metadatatest.AssertPennsieveInstanceIDEqual(t, "60f21224-481c-4e29-a325-c896f184aebe", sample2.PennsieveID)
	metadatatest.AssertPennsieveInstanceIDEqual(t, "60f21224-481c-4e29-a325-c896f184aebe", sample2.GetPennsieveID())

	metadatatest.AssertExternalInstanceIDEqual(t, "09d2a327-be38-403a-884d-a4d1d98b732c", sample2.ID)
	metadatatest.AssertExternalInstanceIDEqual(t, "09d2a327-be38-403a-884d-a4d1d98b732c", sample2.ExternalID())

}
