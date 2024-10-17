package fromcuration

import (
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/mappings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMapProxies(t *testing.T) {
	curationExportPath := "testdata/curation-export.json"

	datasetExport, err := UnmarshalDatasetExport(curationExportPath)
	require.NoError(t, err)

	samples, err := mappings.MapSlice(datasetExport.Samples, ToSample)
	require.NoError(t, err)

	subjects, err := mappings.MapSlice(datasetExport.Subjects, ToSubject)
	require.NoError(t, err)

	proxies, err := MapProxies(samples, subjects, datasetExport.SpecimenDirs.Records, datasetExport.DirStructure)
	require.NoError(t, err)

	assert.Len(t, proxies, 3)
	assert.Contains(t, proxies, metadata.Proxy{
		PackageNodeID: "N:collection:57e23ec5-824f-4690-b615-83dbd5e4d626",
		ProxyKey: metadata.ProxyKey{
			ModelName: metadata.SampleModelName,
			EntityID:  "09d2a327-be38-403a-884d-a4d1d98b732c",
		},
	})

	assert.Contains(t, proxies, metadata.Proxy{
		PackageNodeID: "N:collection:dfca46a1-ceba-44f5-a462-4ec9d61b6a5f",
		ProxyKey: metadata.ProxyKey{
			ModelName: metadata.SampleModelName,
			EntityID:  "09d2a327-be38-403a-884d-a4d1d98b732c",
		},
	})

	assert.Contains(t, proxies, metadata.Proxy{
		ProxyKey: metadata.ProxyKey{
			ModelName: metadata.SubjectModelName,
			EntityID:  "f1027e6e-17cf-45d7-8b57-4c91bfd93fce",
		},
		PackageNodeID: "N:collection:d1011339-b3a0-495d-bf94-76c1a1da1872",
	})
}
