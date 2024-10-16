package fromrecord

import (
	"github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMapProxies_Samples(t *testing.T) {
	inputDirectory := "testdata/input"

	reader, err := client.NewReader(inputDirectory)
	require.NoError(t, err)

	samples, err := MapRecords(reader, metadata.SampleModelName, ToSample)
	require.NoError(t, err)

	proxies, err := MapProxies(reader, metadata.SampleModelName, samples)
	require.NoError(t, err)
	assert.Len(t, proxies, 2)

	assert.Contains(t, proxies, metadata.SavedProxy{
		PennsieveID: "1758ac9d-38db-4c95-b356-04900851b761",
		RecordID:    "cf811b54-bab7-49f1-b239-79c0f6cac29a",
		Proxy: metadata.Proxy{
			EntityID:      "sample-689",
			PackageNodeID: "N:collection:908d58fd-878f-4c98-8b15-b8092628c79e",
		},
	})

	assert.Contains(t, proxies, metadata.SavedProxy{
		PennsieveID: "240951e1-1633-4c26-bfc5-dbbe2dbf17dc",
		RecordID:    "cf811b54-bab7-49f1-b239-79c0f6cac29a",
		Proxy: metadata.Proxy{
			EntityID:      "sample-689",
			PackageNodeID: "N:collection:4e21c210-1880-4ce7-b714-73fc969d66b9",
		},
	})
}

func TestMapProxies_Subjects(t *testing.T) {
	inputDirectory := "testdata/input"

	reader, err := client.NewReader(inputDirectory)
	require.NoError(t, err)

	subjects, err := MapRecords(reader, metadata.SubjectModelName, ToSubject)
	require.NoError(t, err)

	proxies, err := MapProxies(reader, metadata.SubjectModelName, subjects)
	require.NoError(t, err)
	assert.Len(t, proxies, 1)

	assert.Contains(t, proxies, metadata.SavedProxy{
		PennsieveID: "3e078a27-531d-45a4-9ad8-3b166cc86324",
		RecordID:    "a6725f6b-4504-490f-90bc-f21765d0cb07",
		Proxy: metadata.Proxy{
			EntityID:      "dog-123",
			PackageNodeID: "N:collection:45229042-690d-4b4a-9a78-7c8e37a6ff20",
		},
	})

}
