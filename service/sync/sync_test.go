package sync

import (
	"fmt"
	"github.com/google/uuid"
	changesetmodels "github.com/pennsieve/processor-post-metadata/client/models"
	metadataclient "github.com/pennsieve/processor-pre-metadata/client"
	"github.com/pennsieve/ttl-sync-processor/client/metadatatest"
	"github.com/pennsieve/ttl-sync-processor/client/models/metadata"
	"github.com/pennsieve/ttl-sync-processor/service/mappings/fromrecord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestComputeChangeset(t *testing.T) {
	for scenario, test := range map[string]func(tt *testing.T){
		"handle edge cases without panic": emptyChangesets,
		"smoke test":                      smokeTest,
	} {
		t.Run(scenario, func(t *testing.T) {
			ExistingRecordStore = fromrecord.NewRecordIDStore()
			recordIDMap = make(fromrecord.RecordIDMap)

			test(t)
		})
	}
}

func emptyChangesets(t *testing.T) {
	changes, err := ComputeChangeset(emptySchema, &metadata.SavedDatasetMetadata{}, &metadata.DatasetMetadata{})
	require.NoError(t, err)
	require.NotNil(t, changes)
	assert.Empty(t, changes.Models)
	assert.Empty(t, changes.LinkedProperties)
}

func smokeTest(t *testing.T) {
	schemaData := metadataclient.NewSchema(newTestSchemaData().
		WithModel(metadata.ContributorModelName, metadata.ContributorDisplayName).
		WithModel(metadata.SubjectModelName, metadata.SubjectDisplayName).
		WithModel(metadata.SampleModelName, metadata.SampleDisplayName).
		WithLinkedProperty(metadata.SampleSubjectLinkName, metadata.SampleSubjectLinkDisplayName).
		WithProxyRelationshipSchema().
		Build(),
	)

	contributor := metadatatest.NewContributorBuilder().WithNodeID().Build()
	subject := metadatatest.NewSubjectBuilder().Build()
	sample := metadatatest.NewSampleBuilder().Build()
	sampleSubject := metadata.SampleSubject{
		SampleID:  sample.ExternalID(),
		SubjectID: subject.ExternalID(),
	}
	// These are not part of the changeset. But just used to test creating a sample subject
	// link on records that already exist and are not part of the changeset
	existingSubject := metadatatest.NewSavedSubject(metadatatest.NewSubjectBuilder().Build())
	addSavedSubjectRecordID(ExistingRecordStore, existingSubject)
	existingSample := metadatatest.NewSavedSample(metadatatest.NewSampleBuilder().Build())
	addSavedSampleRecordID(ExistingRecordStore, existingSample)

	sampleSubjectOnExisting := metadata.SampleSubject{
		SampleID:  existingSample.ExternalID(),
		SubjectID: existingSubject.ExternalID(),
	}
	proxy := metadata.Proxy{
		ProxyKey: metadata.ProxyKey{
			ModelName:        metadata.SampleModelName,
			TargetExternalID: metadatatest.NewExternalInstanceID(),
		},
		PackageNodeID: fmt.Sprintf("N:collection:%s", uuid.NewString()),
	}

	changes, err := ComputeChangeset(schemaData,
		&metadata.SavedDatasetMetadata{},
		&metadata.DatasetMetadata{
			Contributors:   []metadata.Contributor{contributor},
			Subjects:       []metadata.Subject{subject},
			Samples:        []metadata.Sample{sample},
			SampleSubjects: []metadata.SampleSubject{sampleSubject, sampleSubjectOnExisting},
			Proxies:        []metadata.Proxy{proxy},
		},
	)
	require.NoError(t, err)

	assert.Len(t, changes.Models, 3)
	for _, m := range changes.Models {
		require.NotNil(t, m.ID)
		assert.Len(t, m.Records.Create, 1)
	}

	assert.Len(t, changes.LinkedProperties, 1)
	sampleSubjectChanges := changes.LinkedProperties[0]
	assert.NotNil(t, sampleSubjectChanges.ID)
	assert.Len(t, sampleSubjectChanges.Instances.Create, 2)

	assert.NotNil(t, changes.Proxies)
	assert.False(t, changes.Proxies.CreateProxyRelationshipSchema)
	assert.Len(t, changes.Proxies.RecordChanges, 1)

	assert.Len(t, changes.RecordIDMaps, 2)
	assert.Contains(t, changes.RecordIDMaps, changesetmodels.RecordIDMap{
		ModelName: metadata.SampleModelName,
		ExternalToPennsieve: map[changesetmodels.ExternalInstanceID]changesetmodels.PennsieveInstanceID{
			existingSample.ExternalID(): existingSample.GetPennsieveID(),
		},
	})
	assert.Contains(t, changes.RecordIDMaps, changesetmodels.RecordIDMap{
		ModelName: metadata.SubjectModelName,
		ExternalToPennsieve: map[changesetmodels.ExternalInstanceID]changesetmodels.PennsieveInstanceID{
			existingSubject.ExternalID(): existingSubject.GetPennsieveID(),
		},
	})
}
